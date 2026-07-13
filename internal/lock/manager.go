package lock

import (
	"os"
	"path/filepath"
)

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
	dir := filepath.Join(m.locksDir, repoID)
	if err := os.MkdirAll(dir, 0700); err != nil {
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
