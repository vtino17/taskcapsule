//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func findProjectRoot(t *testing.T) string {
	t.Helper()
	wd, _ := os.Getwd()
	root := wd
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			return root
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("cannot find go.mod (project root)")
		}
		root = parent
	}
	return ""
}

func buildBinary(t *testing.T) string {
	t.Helper()
	projectRoot := findProjectRoot(t)
	binName := "taskcapsule"
	if os.PathSeparator == '\\' {
		binName = "taskcapsule.exe"
	}
	path := filepath.Join(t.TempDir(), binName)
	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", path, "./cmd/taskcapsule")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}
	return path
}

func createGitRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()

	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "TaskCapsule Test"},
		{"git", "branch", "-M", "main"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v\n%s", args, err, out)
		}
	}

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit failed: %v\n%s", err, out)
	}

	taskCapsuleHome := t.TempDir()
	return dir, taskCapsuleHome
}

func taskcapsuleCmd(bin, repoDir, homeDir string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"TASKCAPSULE_HOME="+homeDir,
	)
	return cmd
}

func runTaskcapsule(t *testing.T, bin, repoDir, homeDir string, args ...string) (string, int) {
	t.Helper()
	cmd := taskcapsuleCmd(bin, repoDir, homeDir, args...)
	output, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("command taskcapsule %v failed: %v\n%s", args, err, output)
		}
	}
	return string(output), exitCode
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s failed: %v", path, err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func assertExitCode(t *testing.T, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: exit code %d, want %d", msg, got, want)
	}
}
