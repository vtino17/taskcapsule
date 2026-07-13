package checks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Result struct {
	Command  string
	ExitCode int
	Duration time.Duration
	LogPath  string
	Output   string
}

func Run(worktreePath string, command []string) (*Result, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("command must not be empty")
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = worktreePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	output := stdout.String() + stderr.String()

	return &Result{
		Command:  fmt.Sprintf("%s %s", command[0], command[1:]),
		ExitCode: exitCode,
		Duration: duration,
		Output:   output,
	}, nil
}

func SaveLog(logDir string, result *Result) (string, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format("20060102-150405")
	logFile := filepath.Join(logDir, timestamp+".log")

	content := fmt.Sprintf("Command: %s\nDuration: %.1fs\nExit code: %d\n\n%s",
		result.Command, result.Duration.Seconds(), result.ExitCode, result.Output)

	if err := os.WriteFile(logFile, []byte(content), 0644); err != nil {
		return "", err
	}

	return logFile, nil
}
