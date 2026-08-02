package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/vtino17/taskcapsule/internal/capsule"
)

var repoIDPattern = regexp.MustCompile(`^[a-f0-9]{16}$`)

type Manager struct {
	locksDir string
}

func NewManager(stateBase string) *Manager {
	return &Manager{
		locksDir: filepath.Join(stateBase, "locks"),
	}
}

func (m *Manager) lockPath(repoID, capsuleName string) string {
	dir := filepath.Join(m.locksDir, repoID)
	return filepath.Join(dir, capsuleName+".lock")
}

func (m *Manager) Acquire(repoID, capsuleName, command string) (*Lock, error) {
	if !repoIDPattern.MatchString(repoID) {
		return nil, fmt.Errorf("invalid repository ID")
	}
	if err := capsule.ValidateName(capsuleName); err != nil {
		return nil, err
	}
	dir := filepath.Join(m.locksDir, repoID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return nil, err
	}

	path := m.lockPath(repoID, capsuleName)
	return Acquire(path, command)
}

func (m *Manager) Release(lock *Lock) error {
	if lock == nil {
		return nil
	}
	return lock.Release()
}
