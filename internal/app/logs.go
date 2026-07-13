package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vtino17/taskcapsule/internal/git"
)

func ShowLogs(name string, opts LogOptions) ([]byte, error) {
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	stateBase, _ := getStateDir()
	logDir := filepath.Join(stateBase, "capsules", repoID, name, "logs")

	if opts.ServiceName != "" {
		// Show specific service log
		logFile := filepath.Join(logDir, opts.ServiceName+".log")
		return readTail(logFile, opts.Lines)
	}

	// Show all service logs concatenated
	var result []byte
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("no logs found for capsule %q", name)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		logPath := filepath.Join(logDir, entry.Name())
		data, err := readTail(logPath, opts.Lines)
		if err != nil {
			continue
		}
		if len(result) > 0 {
			result = append(result, '\n')
			result = append(result, "--- "...)
			result = append(result, entry.Name()...)
			result = append(result, " ---\n"...)
		}
		result = append(result, data...)
	}

	return result, nil
}

func readTail(path string, lines int) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Count lines
	lineCount := 0
	for _, b := range data {
		if b == '\n' {
			lineCount++
		}
	}

	if lineCount <= lines {
		return data, nil
	}

	// Find starting position
	skipLines := lineCount - lines
	pos := 0
	for skipLines > 0 && pos < len(data) {
		if data[pos] == '\n' {
			skipLines--
		}
		pos++
	}
	if pos < len(data) {
		return data[pos:], nil
	}
	return data, nil
}
