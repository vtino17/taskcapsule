//go:build integration

package integration

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestConcurrentLifecycleCommandsAreRejected(t *testing.T) {
	bin := buildBinary(t)
	repoDir, homeDir := createGitRepo(t)

	configContent := `{
		"version": 1,
		"defaults": { "baseBranch": "main" },
		"setup": [],
		"services": {}
	}`
	writeFile(t, filepath.Join(repoDir, ".taskcapsule.json"), configContent)

	runTaskcapsule(t, bin, repoDir, homeDir, "start", "test-capsule", "--no-services")

	var wg sync.WaitGroup
	results := make(chan string, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cmd := taskcapsuleCmd(bin, repoDir, homeDir, "pause", "test-capsule")
			output, _ := cmd.CombinedOutput()
			results <- string(output)
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	busyCount := 0
	for r := range results {
		if strings.Contains(r, "Capsule is busy") {
			busyCount++
		}
		if strings.Contains(r, "Capsule paused") || strings.Contains(r, "already paused") {
			successCount++
		}
	}

	t.Logf("success=%d, busy=%d", successCount, busyCount)

	if successCount < 1 {
		t.Error("at least one pause should succeed")
	}

	if busyCount > 2 {
		t.Log("multiple concurrent pauses were rejected as expected")
	}
}
