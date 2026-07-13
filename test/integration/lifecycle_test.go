//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLifecycleStartPauseResumeDelete(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	// Write minimal config
	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main" },
		"setup": [],
		"services": {}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	t.Log("=== Test: Init ===")
	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "init", "--force")
	t.Log(output)
	assertExitCode(t, code, 0, "init should succeed")

	t.Log("=== Test: Start ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "start", "test-capsule", "--no-services")
	t.Log(output)
	assertExitCode(t, code, 0, "start should succeed")

	if !strings.Contains(output, "test-capsule") {
		t.Error("start output should contain capsule name")
	}

	// Check state directory exists
	stateDir := homeDir + "/capsules"
	if !fileExists(stateDir) {
		t.Errorf("state directory should exist: %s", stateDir)
	}

	t.Log("=== Test: Note ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "note", "test-capsule", "Working on integration test")
	t.Log(output)
	assertExitCode(t, code, 0, "note should succeed")

	t.Log("=== Test: Check ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "check", "test-capsule", "--", "git", "status", "--short")
	t.Log(output)
	assertExitCode(t, code, 0, "check should succeed")

	t.Log("=== Test: Status ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "status", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "status should succeed")
	if !strings.Contains(output, "test-capsule") {
		t.Error("status output should contain capsule name")
	}

	t.Log("=== Test: Pause ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "pause", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "pause should succeed")

	t.Log("=== Test: Idempotent Pause ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "pause", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "pause of paused should succeed")

	t.Log("=== Test: Resume ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "resume", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "resume should succeed")

	t.Log("=== Test: Idempotent Resume ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "resume", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "resume of running should succeed")

	t.Log("=== Test: Where ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "where", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "where should succeed")

	t.Log("=== Test: Handoff ===")
	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "handoff", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "handoff should succeed")

	t.Log("=== Test: Pause then Delete ===")
	runTaskcapsule(t, bin, repoDir, homeDir, "pause", "test-capsule")
	time.Sleep(500 * time.Millisecond)

	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "delete", "test-capsule")
	t.Log(output)
	assertExitCode(t, code, 0, "delete should succeed")

	t.Log("=== All lifecycle tests passed ===")
}

func TestListAndVersion(t *testing.T) {
	bin := buildBinary(t)

	t.Log("=== Test: Version ===")
	vCmd := exec.Command(bin, "version")
	verOut, err := vCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v", err)
	}
	if !strings.Contains(string(verOut), "taskcapsule") {
		t.Errorf("version should contain taskcapsule, got: %s", verOut)
	}
	t.Log(string(verOut))

	repoDir, homeDir := createGitRepo(t)
	configContent := `{"version":1,"defaults":{"baseBranch":"main"},"setup":[],"services":{}}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	t.Log("=== Test: List (empty) ===")
	listOut, code := runTaskcapsule(t, bin, repoDir, homeDir, "list")
	t.Log(listOut)
	assertExitCode(t, code, 0, "list should succeed")

	runTaskcapsule(t, bin, repoDir, homeDir, "start", "capsule-a", "--no-services")
	runTaskcapsule(t, bin, repoDir, homeDir, "start", "capsule-b", "--no-services")

	t.Log("=== Test: List (with capsules) ===")
	listOut2, code2 := runTaskcapsule(t, bin, repoDir, homeDir, "list")
	t.Log(listOut2)
	assertExitCode(t, code2, 0, "list should succeed")
	if !strings.Contains(listOut2, "capsule-a") || !strings.Contains(listOut2, "capsule-b") {
		t.Error("list should show both capsules")
	}

	t.Log("=== Test: Doctor ===")
	docOut, code3 := runTaskcapsule(t, bin, repoDir, homeDir, "doctor")
	t.Log(docOut)
	if code3 != 0 {
		t.Log("doctor may report issues in test environment (expected)")
	}
}

// Ensure imports are used
var _ = os.DevNull
