package app

import (
	"fmt"
	"strings"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Where(name string) (*WhereInfo, error) {
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

	info := &WhereInfo{
		Name:     s.Name,
		Status:   s.Status,
		LastNote: s.CurrentNote,
	}

	// Get modified files
	var modifiedFiles []string
	dirty, _ := git.IsDirty(s.WorktreePath)
	if dirty {
		files, _ := git.ChangedFiles(s.WorktreePath)
		modifiedFiles = files
	} else {
		// Get recently modified files (last 24h)
		files, _ := git.ChangedFiles(s.WorktreePath)
		modifiedFiles = files
	}

	// Limit to top 5
	if len(modifiedFiles) > 5 {
		modifiedFiles = modifiedFiles[:5]
	}
	info.ModifiedFiles = modifiedFiles

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
