//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSecretValuesAreNeverPersisted(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main" },
		"setup": [],
		"services": {}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	os.Setenv("TASKCAPSULE_TEST_SECRET", "super-secret-value-12345")
	defer os.Unsetenv("TASKCAPSULE_TEST_SECRET")

	runTaskcapsule(t, bin, repoDir, homeDir, "start", "secret-capsule", "--no-services")
	runTaskcapsule(t, bin, repoDir, homeDir, "note", "secret-capsule", "Working on secret test")
	runTaskcapsule(t, bin, repoDir, homeDir, "handoff", "secret-capsule")
	runTaskcapsule(t, bin, repoDir, homeDir, "where", "secret-capsule")

	// Search for the secret in all state files
	secret := "super-secret-value-12345"
	found := false
	err := filepath.Walk(homeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), secret) {
			found = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if found {
		t.Error("secret value found in state directory - security violation!")
	}

	// Cleanup
	runTaskcapsule(t, bin, repoDir, homeDir, "delete", "secret-capsule", "--force")
}
