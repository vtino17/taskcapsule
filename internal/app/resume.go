package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Resume(name string, opts ResumeOptions) (*ResumeResult, error) {
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

	stateBase, _ := getStateDir()
	cs := state.NewStore(stateBase)

	s, err := cs.Load(repoID, name)
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
	s, err = cs.Load(repoID, name)
	if err != nil {
		return nil, fmt.Errorf("capsule not found: %s", name)
	}

	if _, err := os.Stat(s.WorktreePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree missing for capsule %q", name)
	}

	s.Status = "resuming"
	s.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, s)

	if opts.RunSetup {
		for _, setup := range cfg.Setup {
			cmd := execInDir(setup.Command, s.WorktreePath)
			if output, err := cmd.CombinedOutput(); err != nil {
				setErrorState(s, fmt.Errorf("setup failed on resume: %v\n%s", err, string(output)), cs, repoID, name)
				return nil, fmt.Errorf("setup failed: %v\n%s", err, string(output))
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
		rollbackServices(started, s, cs)
		setErrorState(s, startErr, cs, repoID, name)
		return nil, startErr
	}

	s.Status = "running"
	s.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, s)

	return &ResumeResult{
		Services: serviceInfos,
		LastNote: s.CurrentNote,
	}, nil
}
