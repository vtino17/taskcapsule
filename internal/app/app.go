package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ExitSuccess    = 0
	ExitFailure    = 1
	ExitUsage      = 2
	ExitNotFound   = 3
	ExitUnsafe     = 4
	ExitDependency = 5
	ExitInternal   = 10
)

type CapsuleInfo struct {
	Name    string
	Status  string
	Branch  string
	Updated string
}

type ServiceInfo struct {
	Name    string
	PID     int
	Port    int
	Running bool
}

type CheckInfo struct {
	Command    string
	ExitCode   int
	Duration   float64
	LogPath    string
	FinishedAt string
}

type StatusInfo struct {
	Name           string
	Status         string
	RepositoryRoot string
	WorktreePath   string
	Branch         string
	BaseBranch     string
	Dirty          string
	ChangedFiles   int
	Services       []ServiceInfo
	LastCheck      *CheckInfo
	LastNote       string
}

type WhereInfo struct {
	Name          string
	Status        string
	LastNote      string
	ModifiedFiles []string
	LastCheck     *CheckInfo
}

type StartResult struct {
	Name         string
	Branch       string
	WorktreePath string
	Status       string
	Services     []ServiceInfo
}

type PauseResult struct {
	AlreadyPaused bool
	Services      []ServiceInfo
}

type ResumeResult struct {
	AlreadyRunning bool
	Services       []ServiceInfo
	LastNote       string
}

type CheckResult struct {
	Command  string
	ExitCode int
	Duration float64
	LogPath  string
}

type DoctorResult struct {
	OK      bool
	Message string
}

type LogOptions struct {
	ServiceName string
	Lines       int
	Follow      bool
}

type StartOptions struct {
	BaseBranch string
	Branch     string
	NoServices bool
}

type ResumeOptions struct {
	RunSetup bool
}

func gitExec(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func findGitRoot() (string, error) {
	out, err := gitExec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not a git repository (or git not found)")
	}
	return out, nil
}

func getStateDir() (string, error) {
	if dir := os.Getenv("TASKCAPSULE_HOME"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %v", err)
	}
	return filepath.Join(home, ".taskcapsule"), nil
}
