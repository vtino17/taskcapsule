package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func RunCheck(name string, cmdArgs []string) (*CheckResult, error) {
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	cl, err := acquireCapsuleLock(repoID, name, "check")
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

	worktreePath := s.WorktreePath
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree not found: %s", worktreePath)
	}

	checkDir := filepath.Join(stateBase, "capsules", repoID, name, "checks")
	os.MkdirAll(checkDir, 0755)

	timestamp := time.Now().UTC().Format("20060102-150405")
	logFile := filepath.Join(checkDir, timestamp+".log")
	start := time.Now()

	cmd := execInDir(cmdArgs, worktreePath)
	var outputBuf strings.Builder
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err = cmd.Run()
	duration := time.Since(start).Seconds()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	logContent := fmt.Sprintf("Command: %s\nStarted: %s\nDuration: %.1fs\nExit code: %d\n\n%s",
		strings.Join(cmdArgs, " "), start.Format(time.RFC3339), duration, exitCode, outputBuf.String())
	os.WriteFile(logFile, []byte(logContent), 0644)

	s.LastCheck = &capsule.CheckState{
		Command:    cmdArgs,
		ExitCode:   exitCode,
		StartedAt:  start.Format(time.RFC3339),
		FinishedAt: time.Now().UTC().Format(time.RFC3339),
		LogPath:    logFile,
	}
	s.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, s)

	return &CheckResult{
		Command:  strings.Join(cmdArgs, " "),
		ExitCode: exitCode,
		Duration: duration,
		LogPath:  logFile,
	}, nil
}
