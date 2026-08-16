package state

import (
	"os"
	"testing"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
)

const testRepoID = "0123456789abcdef"

func TestSaveAndLoad(t *testing.T) {
	dir, err := os.MkdirTemp("", "taskcapsule-state-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewStore(dir)
	state := &capsule.State{
		SchemaVersion: 1,
		Name:          "test-capsule",
		Status:        "running",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := store.Save(testRepoID, "test-capsule", state); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load(testRepoID, "test-capsule")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != "test-capsule" {
		t.Errorf("expected test-capsule, got %s", loaded.Name)
	}
	if loaded.Status != "running" {
		t.Errorf("expected running, got %s", loaded.Status)
	}
}

func TestLoadNotFound(t *testing.T) {
	dir, err := os.MkdirTemp("", "taskcapsule-state-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewStore(dir)
	_, err = store.Load(testRepoID, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent capsule")
	}
}

func TestDelete(t *testing.T) {
	dir, err := os.MkdirTemp("", "taskcapsule-state-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewStore(dir)
	state := &capsule.State{
		SchemaVersion: 1,
		Name:          "test",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	store.Save(testRepoID, "test", state)
	if err := store.Delete(testRepoID, "test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Load(testRepoID, "test")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestList(t *testing.T) {
	dir, err := os.MkdirTemp("", "taskcapsule-state-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewStore(dir)
	store.Save(testRepoID, "capsule-a", &capsule.State{Name: "capsule-a", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	store.Save(testRepoID, "capsule-b", &capsule.State{Name: "capsule-b", CreatedAt: time.Now(), UpdatedAt: time.Now()})

	capsules, err := store.List(testRepoID)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(capsules) != 2 {
		t.Errorf("expected 2 capsules, got %d", len(capsules))
	}
}

func TestRejectsUnsafeStorageKeys(t *testing.T) {
	store := NewStore(t.TempDir())
	value := &capsule.State{SchemaVersion: 1, Name: "safe"}

	tests := []struct {
		repoID string
		name   string
	}{
		{repoID: "../outside", name: "safe"},
		{repoID: testRepoID, name: "../../outside"},
		{repoID: testRepoID, name: "service/name"},
	}

	for _, test := range tests {
		if err := store.Save(test.repoID, test.name, value); err == nil {
			t.Fatalf("Save(%q, %q) accepted an unsafe storage key", test.repoID, test.name)
		}
		if _, err := store.Load(test.repoID, test.name); err == nil {
			t.Fatalf("Load(%q, %q) accepted an unsafe storage key", test.repoID, test.name)
		}
		if err := store.Delete(test.repoID, test.name); err == nil {
			t.Fatalf("Delete(%q, %q) accepted an unsafe storage key", test.repoID, test.name)
		}
	}
}
