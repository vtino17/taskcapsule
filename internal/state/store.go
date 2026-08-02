package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/vtino17/taskcapsule/internal/capsule"
)

var repoIDPattern = regexp.MustCompile(`^[a-f0-9]{16}$`)

func validateKey(repoID, name string) error {
	if !repoIDPattern.MatchString(repoID) {
		return fmt.Errorf("invalid repository ID")
	}
	if err := capsule.ValidateName(name); err != nil {
		return fmt.Errorf("invalid capsule name: %w", err)
	}
	return nil
}

type Store struct {
	basePath string
}

func NewStore(basePath string) *Store {
	return &Store{basePath: basePath}
}

func (s *Store) capsuleDir(repoID, name string) string {
	return filepath.Join(s.basePath, "capsules", repoID, name)
}

func (s *Store) stateFilePath(repoID, name string) string {
	return filepath.Join(s.capsuleDir(repoID, name), "state.json")
}

func (s *Store) listCapsulesDir(repoID string) string {
	return filepath.Join(s.basePath, "capsules", repoID)
}

func (s *Store) Save(repoID, name string, state *capsule.State) error {
	if err := validateKey(repoID, name); err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("capsule state must not be nil")
	}
	dir := s.capsuleDir(repoID, name)
	if err := EnsureDir(dir, 0700); err != nil {
		return fmt.Errorf("cannot create state directory: %v", err)
	}

	statePath := s.stateFilePath(repoID, name)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal state: %v", err)
	}

	if err := AtomicWrite(statePath, data, 0600); err != nil {
		return fmt.Errorf("cannot write state: %v", err)
	}

	return nil
}

func (s *Store) Load(repoID, name string) (*capsule.State, error) {
	if err := validateKey(repoID, name); err != nil {
		return nil, err
	}
	path := s.stateFilePath(repoID, name)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("capsule %q not found in repository %s", name, repoID)
		}
		return nil, fmt.Errorf("cannot read state: %v", err)
	}

	var state capsule.State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("invalid state file %s: %v", path, err)
	}

	return &state, nil
}

func (s *Store) Delete(repoID, name string) error {
	if err := validateKey(repoID, name); err != nil {
		return err
	}
	dir := s.capsuleDir(repoID, name)
	return os.RemoveAll(dir)
}

func (s *Store) List(repoID string) ([]*capsule.State, error) {
	if !repoIDPattern.MatchString(repoID) {
		return nil, fmt.Errorf("invalid repository ID")
	}
	capsulesDir := s.listCapsulesDir(repoID)
	return s.readCapsuleDir(capsulesDir)
}

func (s *Store) ListAll() ([]*capsule.State, error) {
	capsulesDir := filepath.Join(s.basePath, "capsules")

	entries, err := os.ReadDir(capsulesDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot list capsules: %v", err)
	}

	var all []*capsule.State
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		repoCapsules, err := s.readCapsuleDir(filepath.Join(capsulesDir, entry.Name()))
		if err != nil {
			continue
		}
		all = append(all, repoCapsules...)
	}

	return all, nil
}

func (s *Store) readCapsuleDir(dir string) ([]*capsule.State, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read capsule directory: %v", err)
	}

	var states []*capsule.State
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		statePath := filepath.Join(dir, entry.Name(), "state.json")
		data, err := os.ReadFile(statePath)
		if err != nil {
			continue
		}

		var state capsule.State
		if err := json.Unmarshal(data, &state); err != nil {
			continue
		}
		states = append(states, &state)
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].UpdatedAt.After(states[j].UpdatedAt)
	})

	return states, nil
}
