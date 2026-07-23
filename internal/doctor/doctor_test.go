package doctor

import (
	"testing"
)

func TestRunCheckGit(t *testing.T) {
	results, err := Run(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Error("expected at least one check result")
	}
}

func TestRunCheckStateDir(t *testing.T) {
	results, err := Run(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, r := range results {
		if r.Message == "State directory writable" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected state directory check")
	}
}

func TestRunCheckWorktreeDir(t *testing.T) {
	results, err := Run(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, r := range results {
		if r.Message == "Worktree directory writable" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected worktree directory check")
	}
}

func TestCheckResultFormat(t *testing.T) {
	r := CheckResult{OK: true, Message: "test ok"}
	if r.Message != "test ok" {
		t.Errorf("unexpected message: %s", r.Message)
	}
	if !r.OK {
		t.Error("expected OK")
	}
}
