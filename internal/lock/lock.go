package lock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	ErrCapsuleBusy = errors.New("capsule is busy")
	ErrStaleLock   = errors.New("stale capsule lock detected")
)

type LockData struct {
	PID       int       `json:"pid"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
}

type Lock struct {
	path string
	pid  int
}

func Acquire(path, command string) (*Lock, error) {
	data := LockData{
		PID:       os.Getpid(),
		Command:   command,
		CreatedAt: time.Now().UTC(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lock data: %v", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return handleExistingLock(path)
		}
		return nil, fmt.Errorf("failed to create lock: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(payload); err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("failed to write lock: %v", err)
	}

	return &Lock{path: path, pid: data.PID}, nil
}

func handleExistingLock(path string) (*Lock, error) {
	existing, err := readLock(path)
	if err != nil {
		return nil, fmt.Errorf("invalid lock file: %v", err)
	}

	if isAlive(existing.PID) {
		return nil, fmt.Errorf("%w: %s\n\nAnother TaskCapsule operation is currently running.\n\nCommand: %s\nPID: %d\nStarted: %s",
			ErrCapsuleBusy, path, existing.Command, existing.PID, existing.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil, fmt.Errorf("%w: %s\n\nThe previous TaskCapsule process is no longer running.\n\nRun:\n  taskcapsule doctor",
		ErrStaleLock, path)
}

func (l *Lock) Release() error {
	if l == nil {
		return nil
	}
	return os.Remove(l.path)
}

func readLock(path string) (*LockData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var l LockData
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// isAlive checks if a process with the given PID is running.
// Platform-specific implementations in isAlive_unix.go and isAlive_windows.go.
func isAlive(pid int) bool {
	return isProcessAlive(pid)
}
