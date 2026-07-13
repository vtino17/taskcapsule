package app

import (
	"fmt"
	"strings"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Status(name string) (*StatusInfo, error) {
	root, err := findGitRoot()
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

	info := &StatusInfo{
		Name:           s.Name,
		Status:         s.Status,
		RepositoryRoot: s.RepositoryRoot,
		WorktreePath:   s.WorktreePath,
		Branch:         s.Branch,
		BaseBranch:     s.BaseBranch,
		LastNote:       s.CurrentNote,
	}

	// Git dirty status
	dirty, err := git.IsDirty(s.WorktreePath)
	if err == nil {
		if dirty {
			info.Dirty = "yes"
		} else {
			info.Dirty = "no"
		}
	} else {
		info.Dirty = "unknown"
	}

	// Changed file count
	if dirty {
		files, err := git.ChangedFiles(s.WorktreePath)
		if err == nil {
			info.ChangedFiles = len(files)
		}
	}

	// Services
	for svcName, svcState := range s.Services {
		running := svcState.Status == "running"
		info.Services = append(info.Services, ServiceInfo{
			Name:    svcName,
			PID:     svcState.PID,
			Port:    svcState.Port,
			Running: running,
		})
	}

	// Last check
	if s.LastCheck != nil {
		info.LastCheck = &CheckInfo{
			Command:    strings.Join(s.LastCheck.Command, " "),
			ExitCode:   s.LastCheck.ExitCode,
			FinishedAt: s.LastCheck.FinishedAt,
		}
	}

	return info, nil
}
