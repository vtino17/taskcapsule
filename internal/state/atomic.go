package state

import (
	"os"
)

// AtomicWrite performs an atomic file write using temp file + rename.
// The temp file is created in the same directory to ensure same-filesystem rename.
func AtomicWrite(path string, data []byte, perm os.FileMode) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

// EnsureDir creates a directory if it doesn't exist.
func EnsureDir(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
