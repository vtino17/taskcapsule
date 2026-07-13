package state

import (
	"os"
	"testing"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
)

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

	if err := store.Save("repo1", "test-capsule", state); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load("repo1", "test-capsule")
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
	_, err = store.Load("repo1", "nonexistent")
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

	store.Save("repo1", "test", state)
	if err := store.Delete("repo1", "test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Load("repo1", "test")
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
	store.Save("repo1", "capsule-a", &capsule.State{Name: "capsule-a", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	store.Save("repo1", "capsule-b", &capsule.State{Name: "capsule-b", CreatedAt: time.Now(), UpdatedAt: time.Now()})

	capsules, err := store.List("repo1")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(capsules) != 2 {
		t.Errorf("expected 2 capsules, got %d", len(capsules))
	}
}
