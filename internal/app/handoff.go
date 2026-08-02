package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/report"
	"github.com/vtino17/taskcapsule/internal/state"
)

func GenerateHandoff(name string) (string, error) {
	if err := validateCapsuleName(name); err != nil {
		return "", err
	}

	root, err := findGitRoot()
	if err != nil {
		return "", err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return "", err
	}

	cl, err := acquireCapsuleLock(repoID, name, "handoff")
	if err != nil {
		return "", err
	}
	defer cl.Release()

	stateBase, err := getStateDir()
	if err != nil {
		return "", err
	}
	cs := state.NewStore(stateBase)

	s, err := loadValidatedCapsule(cs, stateBase, repoID, name)
	if err != nil {
		return "", fmt.Errorf("capsule not found: %s", name)
	}

	var changedFiles []string
	dirty, _ := git.IsDirty(s.WorktreePath)
	if dirty {
		files, _ := git.ChangedFiles(s.WorktreePath)
		changedFiles = files
	}

	var svcNames []string
	for svcName := range s.Services {
		svcNames = append(svcNames, svcName)
	}

	var lastCheckCmd string
	var lastCheckResult string
	var lastCheckExit int
	if s.LastCheck != nil {
		lastCheckCmd = strings.Join(s.LastCheck.Command, " ")
		lastCheckResult = "failed"
		if s.LastCheck.ExitCode == 0 {
			lastCheckResult = "passed"
		}
		lastCheckExit = s.LastCheck.ExitCode
	}

	handoffMD := report.GenerateHandoff(report.HandoffData{
		Name:            s.Name,
		Status:          s.Status,
		Branch:          s.Branch,
		BaseBranch:      s.BaseBranch,
		Dirty:           dirty,
		ChangedFiles:    changedFiles,
		CurrentNote:     s.CurrentNote,
		Services:        svcNames,
		LastCheckCmd:    lastCheckCmd,
		LastCheckResult: lastCheckResult,
		LastCheckExit:   lastCheckExit,
	})

	handoffDir := filepath.Join(filepath.Dir(s.WorktreePath), "handoffs")
	if err := state.EnsureDir(handoffDir, 0700); err != nil {
		return "", fmt.Errorf("cannot create handoff directory: %v", err)
	}
	destPath := filepath.Join(handoffDir, name+".md")

	projectHandoffDir := filepath.Join(root, ".taskcapsule", "handoff")
	if err := state.EnsureDir(projectHandoffDir, 0700); err != nil {
		return "", fmt.Errorf("cannot create project handoff directory: %v", err)
	}
	projectPath := filepath.Join(projectHandoffDir, name+".md")

	if err := os.WriteFile(destPath, []byte(handoffMD), 0600); err != nil {
		return "", fmt.Errorf("cannot write handoff: %v", err)
	}
	if err := os.WriteFile(projectPath, []byte(handoffMD), 0600); err != nil {
		return "", fmt.Errorf("cannot write project handoff: %v", err)
	}

	return destPath, nil
}
