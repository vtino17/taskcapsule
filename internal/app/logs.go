package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vtino17/taskcapsule/internal/git"
)

var defaultLogReader = DefaultLogReader

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
	if lines <= 0 {
		lines = defaultTailLines
	}
	reader := &LogReader{MaxBytes: defaultTailBytes, MaxLines: lines}
	data, err := reader.ReadTail(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err == nil && info.Size() > defaultTailBytes {
		msg := fmt.Sprintf("... (truncated, showing last %d bytes / %d lines) ...\n", defaultTailBytes, lines)
		data = append([]byte(msg), data...)
	}

	return data, nil
}
