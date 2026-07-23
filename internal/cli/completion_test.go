package cli

import (
	"bytes"
	"os"
	"sort"
	"strings"
	"testing"
)

func captureStdoutStderr(f func()) (string, string) {
	oldOut, oldErr := os.Stdout, os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	f()

	wOut.Close()
	wErr.Close()
	var bufOut, bufErr bytes.Buffer
	bufOut.ReadFrom(rOut)
	bufErr.ReadFrom(rErr)
	os.Stdout = oldOut
	os.Stderr = oldErr
	return bufOut.String(), bufErr.String()
}

func expectedCommands() []string {
	names := commandNames()
	sort.Strings(names)
	return names
}

func TestCompletionBash(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"bash"})
		if code != 0 {
			t.Errorf("exit code: got %d, want 0", code)
		}
	})
	if stderr != "" {
		t.Errorf("unexpected stderr: %s", stderr)
	}
	if stdout == "" {
		t.Fatal("empty stdout")
	}
	if !strings.Contains(stdout, "_taskcapsule") {
		t.Error("missing function declaration")
	}
	if !strings.Contains(stdout, "complete -F") {
		t.Error("missing complete command")
	}
	for _, cmd := range expectedCommands() {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("missing command: %s", cmd)
		}
	}
}

func TestCompletionZsh(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"zsh"})
		if code != 0 {
			t.Errorf("exit code: got %d, want 0", code)
		}
	})
	if stderr != "" {
		t.Errorf("unexpected stderr: %s", stderr)
	}
	if stdout == "" {
		t.Fatal("empty stdout")
	}
	if !strings.Contains(stdout, "#compdef taskcapsule") {
		t.Error("missing header")
	}
	for _, cmd := range expectedCommands() {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("missing command: %s", cmd)
		}
	}
}

func TestCompletionFish(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"fish"})
		if code != 0 {
			t.Errorf("exit code: got %d, want 0", code)
		}
	})
	if stderr != "" {
		t.Errorf("unexpected stderr: %s", stderr)
	}
	if stdout == "" {
		t.Fatal("empty stdout")
	}
	if !strings.Contains(stdout, "complete -c taskcapsule") {
		t.Error("missing complete command")
	}
	for _, cmd := range expectedCommands() {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("missing command: %s", cmd)
		}
	}
}

func TestCompletionPowershell(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"powershell"})
		if code != 0 {
			t.Errorf("exit code: got %d, want 0", code)
		}
	})
	if stderr != "" {
		t.Errorf("unexpected stderr: %s", stderr)
	}
	if stdout == "" {
		t.Fatal("empty stdout")
	}
	count := strings.Count(stdout, "Register-ArgumentCompleter")
	if count != 1 {
		t.Errorf("Register-ArgumentCompleter count: got %d, want 1", count)
	}
	for _, cmd := range expectedCommands() {
		if !strings.Contains(stdout, cmd) {
			t.Errorf("missing command: %s", cmd)
		}
	}
}

func TestCompletionMissingShell(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{})
		if code != 2 {
			t.Errorf("exit code: got %d, want 2", code)
		}
	})
	if stdout != "" {
		t.Error("expected empty stdout")
	}
	if !strings.Contains(stderr, "Usage") {
		t.Error("expected usage in stderr")
	}
}

func TestCompletionUnknownShell(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"unknown"})
		if code != 2 {
			t.Errorf("exit code: got %d, want 2", code)
		}
	})
	if stdout != "" {
		t.Error("expected empty stdout")
	}
	if !strings.Contains(stderr, "unknown") {
		t.Error("expected error message in stderr")
	}
}

func TestCompletionExtraArgs(t *testing.T) {
	stdout, stderr := captureStdoutStderr(func() {
		code := handleCompletion([]string{"bash", "extra"})
		if code != 2 {
			t.Errorf("exit code: got %d, want 2", code)
		}
	})
	if stdout != "" {
		t.Error("expected empty stdout")
	}
	if !strings.Contains(stderr, "too many arguments") {
		t.Error("expected error about extra args in stderr")
	}
}

func TestCompletionIncludesCompletionCommand(t *testing.T) {
	names := commandNames()
	found := false
	for _, n := range names {
		if n == "completion" {
			found = true
			break
		}
	}
	if !found {
		t.Error("completion command excluded from commandNames")
	}
}

func TestCommandNamesSorted(t *testing.T) {
	names := commandNames()
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Errorf("not sorted: %s > %s", names[i-1], names[i])
		}
	}
}

func TestCommandNamesNoDuplicates(t *testing.T) {
	names := commandNames()
	seen := make(map[string]bool)
	for _, n := range names {
		if seen[n] {
			t.Errorf("duplicate: %s", n)
		}
		seen[n] = true
	}
}

func TestCompletionDeterministic(t *testing.T) {
	out1, _ := captureStdoutStderr(func() { handleCompletion([]string{"bash"}) })
	out2, _ := captureStdoutStderr(func() { handleCompletion([]string{"bash"}) })
	if out1 != out2 {
		t.Error("bash completion output is not deterministic")
	}
}
