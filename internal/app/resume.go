package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Resume(name string, opts ResumeOptions) (*ResumeResult, error) {
	if err := validateCapsuleName(name); err != nil {
		return nil, err
	}

	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(filepath.Join(root, ".taskcapsule.json"))
	if err != nil {
		return nil, err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	stateBase, err := getStateDir()
	if err != nil {
		return nil, err
	}
	cs := state.NewStore(stateBase)

	s, err := loadValidatedCapsule(cs, stateBase, repoID, name)
	if err != nil {
		return nil, fmt.Errorf("capsule not found: %s", name)
	}

	if s.Status == "running" {
		return &ResumeResult{AlreadyRunning: true}, nil
	}

	if s.Status != "paused" && s.Status != "error" {
		return nil, fmt.Errorf("capsule %q is in state %q; cannot resume (expected paused or error)", name, s.Status)
	}

	// Acquire lock
	cl, err := acquireCapsuleLock(repoID, name, "resume")
	if err != nil {
		return nil, err
	}
	defer cl.Release()

	// Reload state under lock
	s, err = loadValidatedCapsule(cs, stateBase, repoID, name)
	if err != nil {
		return nil, fmt.Errorf("capsule not found: %s", name)
	}

	if _, err := os.Stat(s.WorktreePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree missing for capsule %q", name)
	}

	s.Status = "resuming"
	s.UpdatedAt = time.Now().UTC()
	if err := cs.Save(repoID, name, s); err != nil {
		return nil, fmt.Errorf("failed to persist resuming state: %v", err)
	}

	if opts.RunSetup {
		for _, setup := range cfg.Setup {
			cmd := execInDir(setup.Command, s.WorktreePath)
			if output, err := cmd.CombinedOutput(); err != nil {
				setupErr := fmt.Errorf("setup failed: %v\n%s", err, string(output))
				persistErr := setErrorState(s, setupErr, cs, repoID, name)
				return nil, errors.Join(setupErr, persistErr)
			}
		}
	}

	// Clear old runtime state
	for svcName := range s.Services {
		svc := s.Services[svcName]
		svc.PID = 0
		svc.ProcessGroup = 0
		svc.Status = "stopped"
		s.Services[svcName] = svc
	}

	// Sequential startup with rollback
	started, serviceInfos, startErr := startServices(cfg, stateBase, repoID, name, s.WorktreePath, s, cs)
	if startErr != nil {
		rollbackServices(started, s)
		persistErr := setErrorState(s, startErr, cs, repoID, name)
		return nil, errors.Join(startErr, persistErr)
	}

	s.Status = "running"
	s.UpdatedAt = time.Now().UTC()
	if err := cs.Save(repoID, name, s); err != nil {
		rollbackServices(started, s)
		persistErr := setErrorState(s, err, cs, repoID, name)
		return nil, errors.Join(fmt.Errorf("services started but running state could not be persisted: %w", err), persistErr)
	}

	return &ResumeResult{
		Services: serviceInfos,
		LastNote: s.CurrentNote,
	}, nil
}
