//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteRefusesDirtyWorktree(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	configContent := `{"version":1,"defaults":{"baseBranch":"main"},"setup":[],"services":{}}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	// Start and pause
	runTaskcapsule(t, bin, repoDir, homeDir, "start", "test-capsule", "--no-services")
	runTaskcapsule(t, bin, repoDir, homeDir, "pause", "test-capsule")

	// Get worktree path from status output
	statusOut, _ := runTaskcapsule(t, bin, repoDir, homeDir, "status", "test-capsule")
	t.Log(statusOut)
	worktreePath := ""
	lines := strings.Split(statusOut, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Worktree:") && i+1 < len(lines) {
			worktreePath = strings.TrimSpace(lines[i+1])
		}
	}
	t.Logf("Found worktree path: %q", worktreePath)

	if worktreePath != "" {
		dirtyFile := filepath.Join(worktreePath, "dirty-file.txt")
		os.WriteFile(dirtyFile, []byte("dirty"), 0644)
		t.Logf("Wrote dirty file to %s", dirtyFile)
	}

	// Try to delete (should fail)
	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "delete", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 4, "delete dirty should be refused")

	if !strings.Contains(output, "uncommitted changes") {
		t.Error("delete should mention uncommitted changes")
	}

	// Force delete should succeed
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "delete", "test-capsule", "--force")
	t.Log(output)
	assertExitCode(t, code, 0, "force delete should succeed")
}

func TestDeleteRefusesRunning(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	configContent := `{"version":1,"defaults":{"baseBranch":"main"},"setup":[],"services":{}}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	runTaskcapsule(t, bin, repoDir, homeDir, "start", "test-capsule", "--no-services")

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "delete", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 4, "delete running should be refused")
}
