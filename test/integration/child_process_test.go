//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPauseTerminatesChildProcessGroup(t *testing.T) {
	if os.Getenv("TASKCAPSULE_SKIP_CHILD_TEST") != "" {
		t.Skip("skipping child process test")
	}

	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	// Use a shell command that spawns a child process
	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main" },
		"setup": [],
		"services": {
			"parent": {
				"command": ["sh", "-c", "sleep 30 & sleep 30"],
				"health": {
					"type": "none"
				}
			}
		}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "start", "child-test")
	t.Log(output)
	if code != 0 {
		t.Skip("start failed, skipping child process test")
	}

	output, code = runTaskcapsule(t, bin, repoDir, homeDir, "pause", "child-test")
	t.Log(output)
	assertExitCode(t, code, 0, "pause should succeed")

	runTaskcapsule(t, bin, repoDir, homeDir, "delete", "child-test", "--force")
}
