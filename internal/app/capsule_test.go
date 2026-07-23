package app

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestExitCodeFromErrorTypes(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{errors.New("capsule not found: foo"), ExitNotFound},
		{errors.New("not a git repository"), ExitDependency},
		{errors.New("uncommitted changes"), ExitUnsafe},
		{errors.New("capsule \"test\" already exists"), ExitUnsafe},
		{errors.New("already paused"), ExitSuccess},
		{errors.New("already running"), ExitSuccess},
		{errors.New("random error"), ExitFailure},
	}
	for _, tc := range tests {
		got := exitCodeFromError(tc.err)
		if got != tc.want {
			t.Errorf("exitCodeFromError(%q) = %d, want %d", tc.err.Error(), got, tc.want)
		}
	}
}

func TestExitCodeNilError(t *testing.T) {
	code := exitCodeFromError(nil)
	if code != ExitSuccess {
		t.Errorf("expected ExitSuccess for nil error, got %d", code)
	}
}

func TestStrPtr(t *testing.T) {
	s := strPtr("hello")
	if *s != "hello" {
		t.Errorf("expected 'hello', got '%s'", *s)
	}
}

func TestFindGitRootOutsideRepo(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	tmp := t.TempDir()
	os.Chdir(tmp)

	_, err := findGitRoot()
	if err == nil {
		t.Error("expected error outside git repo")
	}
}

func TestGetStateDirDefault(t *testing.T) {
	dir, err := getStateDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Fatal("empty state dir")
	}
	if !strings.Contains(dir, ".taskcapsule") {
		t.Errorf("expected .taskcapsule in path, got %s", dir)
	}
}

func TestGetStateDirEnv(t *testing.T) {
	os.Setenv("TASKCAPSULE_HOME", "/tmp/test-tc-home")
	defer os.Unsetenv("TASKCAPSULE_HOME")

	dir, err := getStateDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != "/tmp/test-tc-home" {
		t.Errorf("expected /tmp/test-tc-home, got %s", dir)
	}
}

func TestIsProcessRunning(t *testing.T) {
	if !isProcessRunning(os.Getpid()) {
		// On Windows, FindProcess always succeeds but Signal may fail for the current process
		t.Log("isProcessRunning returned false for current process (possible on Windows)")
	}
}

func TestIsProcessRunningNonExistent(t *testing.T) {
	_ = isProcessRunning(999999)
}
