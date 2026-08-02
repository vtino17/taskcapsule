package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "taskcapsule-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, cmd := range cmds {
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git setup failed: %v\n%s", err, out)
		}
	}

	// Initial commit
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test"), 0644); err != nil {
		t.Fatal(err)
	}
	c := exec.Command("git", "add", ".")
	c.Dir = dir
	c.Run()
	c = exec.Command("git", "commit", "-m", "initial")
	c.Dir = dir
	c.Run()

	return dir
}

func TestRoot(t *testing.T) {
	dir := setupTestRepo(t)

	root, err := execGitInDir(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want, err := canonicalRoot(dir)
	if err != nil {
		t.Fatalf("canonicalRoot: %v", err)
	}
	if root != want {
		t.Errorf("expected %s, got %s", want, root)
	}
}

func TestRepoIDCanonicalizesSymlinkAliases(t *testing.T) {
	dir := setupTestRepo(t)
	alias := filepath.Join(t.TempDir(), "repo-alias")
	if err := os.Symlink(dir, alias); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}

	directID, err := RepoID(dir)
	if err != nil {
		t.Fatalf("RepoID(direct): %v", err)
	}
	aliasID, err := RepoID(alias)
	if err != nil {
		t.Fatalf("RepoID(alias): %v", err)
	}
	if directID != aliasID {
		t.Fatalf("expected canonical aliases to share an ID, got %s and %s", directID, aliasID)
	}
}

func TestIsDirty(t *testing.T) {
	dir := setupTestRepo(t)

	dirty, err := execGitInDir(dir, "status", "--porcelain")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(dirty) != "" {
		t.Error("expected clean repo")
	}
}

func TestBranchExists(t *testing.T) {
	dir := setupTestRepo(t)

	// Check whichever default branch exists (master or main)
	for _, name := range []string{"main", "master"} {
		exists, err := BranchExists(name, dir)
		if err != nil {
			t.Fatal(err)
		}
		if exists {
			return // found default branch
		}
	}
	t.Error("expected a default branch (main or master) to exist after git init")

	exists, err := BranchExists("nonexistent", dir)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("expected nonexistent branch to not exist")
	}
}
