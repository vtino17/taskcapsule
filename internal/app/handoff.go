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

	stateBase, _ := getStateDir()
	cs := state.NewStore(stateBase)

	s, err := cs.Load(repoID, name)
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
	os.MkdirAll(handoffDir, 0755)
	destPath := filepath.Join(handoffDir, name+".md")

	projectHandoffDir := filepath.Join(root, ".taskcapsule", "handoff")
	os.MkdirAll(projectHandoffDir, 0755)
	projectPath := filepath.Join(projectHandoffDir, name+".md")

	if err := os.WriteFile(destPath, []byte(handoffMD), 0644); err != nil {
		return "", fmt.Errorf("cannot write handoff: %v", err)
	}
	os.WriteFile(projectPath, []byte(handoffMD), 0644)

	return destPath, nil
}
