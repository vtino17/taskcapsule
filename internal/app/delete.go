package app

import (
	"fmt"
	"os"
	"time"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/process"
	"github.com/vtino17/taskcapsule/internal/state"
)

func DeleteCapsule(name string, force bool) error {
	root, err := findGitRoot()
	if err != nil {
		return err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return err
	}

	cl, err := acquireCapsuleLock(repoID, name, "delete")
	if err != nil {
		return err
	}
	defer cl.Release()

	stateBase, _ := getStateDir()
	cs := state.NewStore(stateBase)

	s, err := cs.Load(repoID, name)
	if err != nil {
		return fmt.Errorf("capsule not found: %s", name)
	}

	if s.Status == "running" && !force {
		return fmt.Errorf("cannot delete running capsule: %s\nPause it first: taskcapsule pause %s", name, name)
	}

	if !force {
		dirty, err := git.IsDirty(s.WorktreePath)
		if err == nil && dirty {
			return fmt.Errorf("cannot delete capsule: worktree has uncommitted changes\nUse: taskcapsule delete %s --force", name)
		}
	}

	s.Status = "deleting"
	s.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, s)

	if s.StatusPrevious() == "running" || force {
		for _, svcState := range s.Services {
			if svcState.ProcessGroup > 0 {
				process.StopProcessGroup(svcState.ProcessGroup, 5)
			} else if svcState.PID > 0 {
				process.StopProcess(svcState.PID, 5)
			}
		}
	}

	if err := git.RemoveWorktree(s.WorktreePath); err != nil && !force {
		return fmt.Errorf("failed to remove worktree: %v", err)
	}

	if err := cs.Delete(repoID, name); err != nil && !force {
		return fmt.Errorf("failed to remove state: %v", err)
	}

	// Clean up any leftover worktree directory
	os.RemoveAll(s.WorktreePath)

	return nil
}
