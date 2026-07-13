package app

import (
	"fmt"
	"time"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/process"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Pause(name string) (*PauseResult, error) {
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	cl, err := acquireCapsuleLock(repoID, name, "pause")
	if err != nil {
		return nil, err
	}
	defer cl.Release()

	stateBase, _ := getStateDir()
	cs := state.NewStore(stateBase)

	s, err := cs.Load(repoID, name)
	if err != nil {
		return nil, fmt.Errorf("capsule not found: %s", name)
	}

	if s.Status == "paused" {
		return &PauseResult{AlreadyPaused: true}, nil
	}

	if s.Status != "running" {
		return nil, fmt.Errorf("capsule %q is in state %q; cannot pause", name, s.Status)
	}

	s.Status = "pausing"
	s.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, s)

	var svcInfos []ServiceInfo

	for svcName, svcState := range s.Services {
		svcInfos = append(svcInfos, ServiceInfo{
			Name: svcName,
			PID:  svcState.PID,
		})

		if svcState.Status == "running" {
			if svcState.ProcessGroup > 0 {
				process.StopProcessGroup(svcState.ProcessGroup, 5)
			} else if svcState.PID > 0 {
				process.StopProcess(svcState.PID, 5)
			}
		}

		svcState.Status = "stopped"
		svcState.PID = 0
		svcState.ProcessGroup = 0
		s.Services[svcName] = svcState
	}

	now := time.Now().UTC()
	s.Status = "paused"
	s.UpdatedAt = now
	s.LastPausedAt = &now
	cs.Save(repoID, name, s)

	return &PauseResult{Services: svcInfos}, nil
}
