package lock

import (
	"os"
	"testing"
)

func TestAcquireRelease(t *testing.T) {
	path := t.TempDir() + "/test.lock"

	l, err := Acquire(path, "test")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("lock file should exist after acquire")
	}

	if err := l.Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("lock file should be removed after release")
	}
}

func TestDuplicateAcquire(t *testing.T) {
	path := t.TempDir() + "/test.lock"

	l1, err := Acquire(path, "cmd1")
	if err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l1.Release()

	_, err = Acquire(path, "cmd2")
	if err == nil {
		t.Fatal("expected error for duplicate acquire")
	}
}

func TestStaleLock(t *testing.T) {
	path := t.TempDir() + "/test.lock"

	lockData := `{"pid": 999999999, "command": "old", "createdAt": "2026-01-01T00:00:00Z"}`
	if err := os.WriteFile(path, []byte(lockData), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Acquire(path, "new")
	if err == nil {
		t.Fatal("expected error for stale lock on same path")
	}
}

func TestLockAfterRelease(t *testing.T) {
	path := t.TempDir() + "/test.lock"

	l1, err := Acquire(path, "cmd1")
	if err != nil {
		t.Fatal(err)
	}
	l1.Release()

	l2, err := Acquire(path, "cmd2")
	if err != nil {
		t.Fatalf("second acquire should succeed after release: %v", err)
	}
	defer l2.Release()
}

func TestSeparateLocks(t *testing.T) {
	dir := t.TempDir()

	l1, err := Acquire(dir+"/a.lock", "cmd-a")
	if err != nil {
		t.Fatal(err)
	}
	defer l1.Release()

	l2, err := Acquire(dir+"/b.lock", "cmd-b")
	if err != nil {
		t.Fatal(err)
	}
	defer l2.Release()
}

func TestReleaseNilLock(t *testing.T) {
	var l *Lock = nil
	if err := l.Release(); err != nil {
		t.Fatalf("release nil lock should not error: %v", err)
	}
}
