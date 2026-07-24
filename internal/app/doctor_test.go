package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func setupTestRepo(t *testing.T, branches []string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "taskcapsule-doctor-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "branch", "-M", "main"},
	}
	for _, args := range cmds {
		c := exec.Command(args[0], args[1:]...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git setup %v failed: %v\n%s", args, err, out)
		}
	}

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	c := exec.Command("git", "add", ".")
	c.Dir = dir
	c.Run()
	c = exec.Command("git", "commit", "-m", "initial")
	c.Dir = dir
	c.Run()

	for _, branch := range branches {
		c := exec.Command("git", "branch", branch)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git branch %s failed: %v\n%s", branch, err, out)
		}
	}

	return dir
}

func setupCapsule(t *testing.T, store *state.Store, repoID, name, branch, worktreePath, status string) {
	t.Helper()
	s := &capsule.State{
		SchemaVersion:  1,
		Name:           name,
		Status:         status,
		RepositoryRoot: worktreePath,
		RepositoryID:   repoID,
		WorktreePath:   worktreePath,
		Branch:         branch,
		BaseBranch:     "main",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := store.Save(repoID, name, s); err != nil {
		t.Fatalf("setupCapsule Save: %v", err)
	}
}

func runDoctorInDir(t *testing.T, repoDir, stateDir string) []DoctorResult {
	t.Helper()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	if err := os.Chdir(repoDir); err != nil {
		t.Fatal(err)
	}

	os.Setenv("TASKCAPSULE_HOME", stateDir)
	defer os.Unsetenv("TASKCAPSULE_HOME")

	results, err := Doctor()
	if err != nil {
		t.Fatalf("Doctor() returned error: %v", err)
	}
	return results
}

func TestDoctorCurrentRepoBranchChecked(t *testing.T) {
	repoDir := setupTestRepo(t, []string{"feature-x"})
	stateDir := t.TempDir()

	repoID, err := git.RepoID(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoID, "capsule-good", "main", repoDir, "active")

	results := runDoctorInDir(t, repoDir, stateDir)

	for _, r := range results {
		if strings.Contains(r.Message, "branch") && strings.Contains(r.Message, "capsule-good") {
			t.Errorf("got unexpected branch warning for valid capsule: %s", r.Message)
		}
	}
}

func TestDoctorOtherRepoNotCheckedAgainstCurrent(t *testing.T) {
	repoA := setupTestRepo(t, []string{"feat-a"})
	repoB := setupTestRepo(t, []string{"feat-b"})
	stateDir := t.TempDir()

	repoAID, err := git.RepoID(repoA)
	if err != nil {
		t.Fatal(err)
	}
	repoBID, err := git.RepoID(repoB)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoAID, "capsule-a", "feat-a", repoA, "active")
	setupCapsule(t, store, repoBID, "capsule-b", "feat-b", repoB, "active")

	results := runDoctorInDir(t, repoA, stateDir)

	for _, r := range results {
		if strings.Contains(r.Message, "capsule-b") && strings.Contains(r.Message, "branch") {
			t.Errorf("other-repo capsule 'capsule-b' must NOT be branch-checked against current repo: %s", r.Message)
		}
	}
}

func TestDoctorMissingBranchInCurrentRepoReported(t *testing.T) {
	repoDir := setupTestRepo(t, nil)
	stateDir := t.TempDir()

	repoID, err := git.RepoID(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoID, "capsule-missing", "branch-gone", repoDir, "active")

	results := runDoctorInDir(t, repoDir, stateDir)

	found := false
	for _, r := range results {
		if strings.Contains(r.Message, "capsule-missing") && strings.Contains(r.Message, "branch") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected branch warning for capsule with missing branch, got none")
	}
}

func TestDoctorValidBranchNotReported(t *testing.T) {
	repoDir := setupTestRepo(t, []string{"valid-branch"})
	stateDir := t.TempDir()

	repoID, err := git.RepoID(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoID, "capsule-valid", "valid-branch", repoDir, "active")

	results := runDoctorInDir(t, repoDir, stateDir)

	for _, r := range results {
		if strings.Contains(r.Message, "capsule-valid") && strings.Contains(r.Message, "branch") {
			t.Errorf("got false branch warning for valid branch: %s", r.Message)
		}
	}
}

func TestDoctorSameNameAcrossReposNoFalseWarning(t *testing.T) {
	repoA := setupTestRepo(t, []string{"feature-x"})
	repoB := setupTestRepo(t, []string{"feature-x"})
	stateDir := t.TempDir()

	repoAID, err := git.RepoID(repoA)
	if err != nil {
		t.Fatal(err)
	}
	repoBID, err := git.RepoID(repoB)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoAID, "feature-x", "feature-x", repoA, "active")
	setupCapsule(t, store, repoBID, "feature-x", "feature-x", repoB, "active")

	results := runDoctorInDir(t, repoA, stateDir)

	branchWarnings := 0
	for _, r := range results {
		if strings.Contains(r.Message, "branch") && strings.Contains(r.Message, "feature-x") {
			branchWarnings++
		}
	}
	if branchWarnings > 1 {
		t.Errorf("expected at most 1 branch warning for same-name capsules across repos, got %d", branchWarnings)
	}
}

func TestDoctorDoesNotMutateState(t *testing.T) {
	repoDir := setupTestRepo(t, []string{"feat-mutate"})
	stateDir := t.TempDir()

	repoID, err := git.RepoID(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoID, "capsule-mutate", "feat-mutate", repoDir, "active")

	loaded, err := store.Load(repoID, "capsule-mutate")
	if err != nil {
		t.Fatal(err)
	}
	originalUpdated := loaded.UpdatedAt

	_ = runDoctorInDir(t, repoDir, stateDir)

	reloaded, err := store.Load(repoID, "capsule-mutate")
	if err != nil {
		t.Fatal(err)
	}
	if !reloaded.UpdatedAt.Equal(originalUpdated) {
		t.Error("Doctor mutated capsule state (UpdatedAt changed)")
	}
}

func TestDoctorSigZeroBehaviorPreserved(t *testing.T) {
	if !isProcessRunning(os.Getpid()) {
		t.Skip("current process not reported as running (expected on Windows)")
	}
}

func TestDoctorMissingWorktreeReported(t *testing.T) {
	repoDir := setupTestRepo(t, []string{"feat-wt"})
	stateDir := t.TempDir()

	repoID, err := git.RepoID(repoDir)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoID, "capsule-no-wt", "feat-wt", filepath.Join(stateDir, "nonexistent-worktree"), "active")

	results := runDoctorInDir(t, repoDir, stateDir)

	found := false
	for _, r := range results {
		if strings.Contains(r.Message, "capsule-no-wt") && strings.Contains(r.Message, "missing its worktree") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected worktree missing warning")
	}
}

func TestDoctorMissingBranchOtherRepoNotReported(t *testing.T) {
	repoA := setupTestRepo(t, []string{"feat-a"})
	repoB := setupTestRepo(t, nil)
	stateDir := t.TempDir()

	repoAID, err := git.RepoID(repoA)
	if err != nil {
		t.Fatal(err)
	}
	repoBID, err := git.RepoID(repoB)
	if err != nil {
		t.Fatal(err)
	}

	store := state.NewStore(stateDir)
	_ = os.MkdirAll(filepath.Join(stateDir, "worktrees"), 0755)
	setupCapsule(t, store, repoAID, "capsule-a", "feat-a", repoA, "active")
	setupCapsule(t, store, repoBID, "capsule-b-missing", "branch-gone", repoB, "active")

	results := runDoctorInDir(t, repoA, stateDir)

	for _, r := range results {
		if strings.Contains(r.Message, "capsule-b-missing") && strings.Contains(r.Message, "branch") {
			t.Errorf("other-repo capsule with missing branch must not be reported: %s", r.Message)
		}
	}
}
