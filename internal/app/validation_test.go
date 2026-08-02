package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateManagedWorktreePath(t *testing.T) {
	stateBase := t.TempDir()
	managed := filepath.Join(stateBase, "worktrees", "repo", "capsule")
	if err := validateManagedWorktreePath(stateBase, managed); err != nil {
		t.Fatalf("managed path rejected: %v", err)
	}

	unsafe := []string{
		stateBase,
		filepath.Join(stateBase, "worktrees"),
		filepath.Dir(stateBase),
		filepath.Join(stateBase, "..", "outside"),
	}
	for _, path := range unsafe {
		if err := validateManagedWorktreePath(stateBase, path); err == nil {
			t.Errorf("unsafe path %q was accepted", path)
		}
	}
}

func TestValidateManagedWorktreePathRejectsSymlinkEscape(t *testing.T) {
	stateBase := t.TempDir()
	managedRoot := filepath.Join(stateBase, "worktrees")
	outside := t.TempDir()
	if err := os.MkdirAll(managedRoot, 0700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(managedRoot, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := validateManagedWorktreePath(stateBase, link); err == nil {
		t.Fatal("managed worktree symlink escape was accepted")
	}
}

func TestValidateCapsuleNameRejectsTraversal(t *testing.T) {
	for _, name := range []string{"../escape", "../../escape", "service/name", `service\\name`} {
		if err := validateCapsuleName(name); err == nil {
			t.Errorf("unsafe capsule name %q was accepted", name)
		}
	}
}
