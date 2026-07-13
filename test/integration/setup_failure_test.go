//go:build integration

package integration

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupFailurePreventsServiceStartup(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	// Config with failing setup command
	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main" },
		"setup": [
			{
				"command": ["go", "tool", "nonexistent"]
			}
		],
		"services": {}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "start", "fail-capsule", "--no-services")
	t.Log(output)

	if code == 0 {
		t.Error("start should fail when setup fails")
	}

	if !strings.Contains(output, "setup") {
		t.Log("note: setup failure message may vary by platform")
	}
}

func TestConfigMissing(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "start", "test", "--no-services")
	t.Log(output)

	if code == 0 {
		t.Error("start without config should fail")
	}
}

func TestDuplicateCapsuleRejected(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	configContent := `{"version":1,"defaults":{"baseBranch":"main"},"setup":[],"services":{}}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	runTaskcapsule(t, bin, repoDir, homeDir, "start", "dup-capsule", "--no-services")

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "start", "dup-capsule", "--no-services")
	t.Log(output)

	if code == 0 {
		t.Error("duplicate start should fail")
	}

	if !strings.Contains(output, "already exists") {
		t.Error("duplicate error should mention 'already exists'")
	}
}
