package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/state"
)

func validateCapsuleName(name string) error {
	if err := capsule.ValidateName(name); err != nil {
		return fmt.Errorf("invalid capsule name: %w", err)
	}
	return nil
}

func validateManagedWorktreePath(stateBase, worktreePath string) error {
	if worktreePath == "" {
		return fmt.Errorf("capsule state has an empty worktree path")
	}

	managedRoot, err := filepath.Abs(filepath.Join(stateBase, "worktrees"))
	if err != nil {
		return fmt.Errorf("resolve managed worktree root: %w", err)
	}
	path, err := filepath.Abs(worktreePath)
	if err != nil {
		return fmt.Errorf("resolve capsule worktree path: %w", err)
	}

	rel, err := filepath.Rel(managedRoot, path)
	if err != nil {
		return fmt.Errorf("compare capsule worktree path with managed root: %w", err)
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("unsafe capsule worktree path %q: must be below %q", worktreePath, managedRoot)
	}

	resolvedRoot, rootErr := filepath.EvalSymlinks(managedRoot)
	resolvedPath, pathErr := filepath.EvalSymlinks(path)
	if rootErr == nil && pathErr == nil {
		resolvedRel, err := filepath.Rel(resolvedRoot, resolvedPath)
		if err != nil {
			return fmt.Errorf("compare resolved capsule worktree path with managed root: %w", err)
		}
		if resolvedRel == "." || resolvedRel == ".." || strings.HasPrefix(resolvedRel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("unsafe resolved capsule worktree path %q: must be below %q", resolvedPath, resolvedRoot)
		}
	}

	return nil
}

func loadValidatedCapsule(cs *state.Store, stateBase, repoID, name string) (*capsule.State, error) {
	s, err := cs.Load(repoID, name)
	if err != nil {
		return nil, err
	}
	if s.SchemaVersion != 1 {
		return nil, fmt.Errorf("unsupported capsule state schema version %d", s.SchemaVersion)
	}
	if s.Name != name {
		return nil, fmt.Errorf("capsule state name mismatch: expected %q, got %q", name, s.Name)
	}
	if s.RepositoryID != repoID {
		return nil, fmt.Errorf("capsule state repository mismatch")
	}
	if err := validateManagedWorktreePath(stateBase, s.WorktreePath); err != nil {
		return nil, err
	}
	return s, nil
}
