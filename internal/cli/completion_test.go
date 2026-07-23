package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old
	return buf.String()
}

func TestCompletionBash(t *testing.T) {
	out := captureStdout(func() {
		handleCompletion([]string{"bash"})
	})
	if !strings.Contains(out, "_taskcapsule") {
		t.Error("bash completion missing function")
	}
}

func TestCompletionZsh(t *testing.T) {
	out := captureStdout(func() {
		handleCompletion([]string{"zsh"})
	})
	if !strings.Contains(out, "#compdef taskcapsule") {
		t.Error("zsh completion missing header")
	}
}

func TestCompletionFish(t *testing.T) {
	out := captureStdout(func() {
		handleCompletion([]string{"fish"})
	})
	if !strings.Contains(out, "complete -c taskcapsule") {
		t.Error("fish completion missing complete command")
	}
}

func TestCompletionPowershell(t *testing.T) {
	out := captureStdout(func() {
		handleCompletion([]string{"powershell"})
	})
	if !strings.Contains(out, "Register-ArgumentCompleter") {
		t.Error("powershell completion missing Register-ArgumentCompleter")
	}
}

func TestCompletionUnknownShell(t *testing.T) {
	exit := handleCompletion([]string{"unknown"})
	if exit != 2 {
		t.Errorf("expected exit 2, got %d", exit)
	}
}

func TestCompletionMissingShell(t *testing.T) {
	exit := handleCompletion([]string{})
	if exit != 2 {
		t.Errorf("expected exit 2, got %d", exit)
	}
}

func TestCommandNames(t *testing.T) {
	names := commandNames()
	if len(names) == 0 {
		t.Fatal("expected at least one command")
	}
	if names[0] == "completion" {
		t.Error("completion should not include itself")
	}
	for _, n := range names {
		if n == "" {
			t.Error("empty command name")
		}
	}
}

func TestCommands(t *testing.T) {
	// Verify the registered commands list is correct
	names := commandNames()
	expected := []string{"check", "delete", "doctor", "handoff", "init", "list", "logs", "note", "pause", "resume", "start", "status", "version", "where"}
	for _, exp := range expected {
		found := false
		for _, n := range names {
			if n == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing command: %s", exp)
		}
	}
}
