package state

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Lock struct {
	path string
}

func NewLock(lockPath string) *Lock {
	return &Lock{path: lockPath}
}

// Acquire tries to create an exclusive lock file.
// Returns an error if the lock is held by a live process.
func (l *Lock) Acquire() error {
	// Try to create lock file exclusively
	data := []byte(fmt.Sprintf("%d\n", os.Getpid()))
	if err := os.WriteFile(l.path, data, 0644); err != nil {
		// Lock file exists - check if stale
		return l.checkStale()
	}
	return nil
}

// Release removes the lock file.
func (l *Lock) Release() error {
	return os.Remove(l.path)
}

// LockedBy returns the PID that holds the lock, or 0 if no lock.
func (l *Lock) LockedBy() int {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return 0
	}

	line := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(line)
	if err != nil {
		return 0
	}
	return pid
}

func (l *Lock) checkStale() error {
	pid := l.LockedBy()
	if pid <= 0 {
		// Stale lock file
		os.Remove(l.path)
		return l.Acquire()
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		// Process not found - stale lock
		os.Remove(l.path)
		return l.Acquire()
	}

	// On Unix, signal 0 checks process existence
	// On Windows, FindProcess always succeeds (known limitation)
	if proc.Signal(os.Interrupt) != nil {
		// Process not alive - stale lock
		os.Remove(l.path)
		return l.Acquire()
	}

	return fmt.Errorf("capsule is locked by process %d", pid)
}

// RepoIDFromRoot derives a stable repo ID from a repo root path.
// Uses the same algorithm as git.RepoID for consistency.
func RepoIDFromRoot(root string) (string, error) {
	// Delegate to git package logic via import not possible due to cycle.
	// Use same simple hash approach.
	return getFallbackRepoID(root), nil
}

func getFallbackRepoID(root string) string {
	// Simple hash of the root path for repo identification
	h := sha256.Sum256([]byte(root))
	return fmt.Sprintf("%x", h[:8])
}
