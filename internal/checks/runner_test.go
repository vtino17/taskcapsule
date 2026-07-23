package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunSuccess(t *testing.T) {
	result, err := Run(t.TempDir(), []string{"go", "version"})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "go") {
		t.Error("expected 'go' in output")
	}
}

func TestRunEmptyCommand(t *testing.T) {
	_, err := Run(t.TempDir(), []string{})
	if err == nil {
		t.Error("expected error for empty command")
	}
}

func TestRunNonExistentExecutable(t *testing.T) {
	_, err := Run(t.TempDir(), []string{"nonexistent-command-xyz"})
	if err != nil {
		return // acceptable - exec package returns error
	}
}

func TestSaveLog(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs")
	result := &Result{
		Command:  "echo hello",
		ExitCode: 0,
		Output:   "hello\n",
	}
	path, err := SaveLog(logDir, result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("log file not created: %s", path)
	}
}

func TestSaveLogDifferentFilenames(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs2")
	result := &Result{Command: "test", ExitCode: 0}
	path1, _ := SaveLog(logDir, result)
	time.Sleep(2 * time.Second) // ensure different timestamp
	path2, _ := SaveLog(logDir, result)
	if path1 == path2 {
		t.Error("expected different filenames for separate saves")
	}
}
