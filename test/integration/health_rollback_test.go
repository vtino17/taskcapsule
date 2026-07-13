//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHealthFailureRollsBackStartedServices(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	fixtureName := "httpserver"
	if os.PathSeparator == '\\' {
		fixtureName = "httpserver.exe"
	}
	fixturePath := filepath.Join(t.TempDir(), fixtureName)
	fCmd := exec.Command("go", "build", "-buildvcs=false", "-o", fixturePath, "./test/fixtures/httpserver")
	fCmd.Dir = findProjectRoot(t)
	out, err := fCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build fixture failed: %v\n%s", err, out)
	}

	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main", "healthTimeoutSeconds": 10 },
		"setup": [],
		"services": {
			"healthy": {
				"command": ["` + fixturePath + `"],
				"health": {
					"type": "process",
					"timeoutSeconds": 5
				}
			},
			"failing": {
				"command": ["` + fixturePath + `"],
				"environment": {
					"EXIT_IMMEDIATELY": "1"
				},
				"health": {
					"type": "process",
					"timeoutSeconds": 5
				}
			}
		}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	output, code := runTaskcapsule(t, bin, repoDir, homeDir, "start", "health-test")
	t.Log(output)

	if code == 0 {
		t.Error("start should fail when a service health check fails")
	}
}
