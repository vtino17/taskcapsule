package config

import (
	"os"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	content := `{
		"version": 1,
		"defaults": {
			"baseBranch": "main",
			"branchPrefix": "task/",
			"gracefulShutdownSeconds": 5,
			"healthTimeoutSeconds": 30
		},
		"services": {
			"api": {
				"command": ["go", "run", "./cmd/api"]
			}
		}
	}`

	f, err := os.CreateTemp("", "taskcapsule-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if _, err := f.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if cfg.Defaults.BaseBranch != "main" {
		t.Errorf("expected baseBranch main, got %s", cfg.Defaults.BaseBranch)
	}
}

func TestLoadUnknownSchemaVersion(t *testing.T) {
	content := `{"version": 999}`

	f, err := os.CreateTemp("", "taskcapsule-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.Write([]byte(content))
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Fatal("expected error for unknown schema version")
	}
}

func TestLoadDuplicateService(t *testing.T) {
	content := `{
		"version": 1,
		"services": {
			"api": { "command": ["go", "run"] },
			"api": { "command": ["node", "server.js"] }
		}
	}`

	// Note: JSON will override duplicate keys during unmarshal
	// so this test validates the map doesn't have issues
	f, err := os.CreateTemp("", "taskcapsule-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.Write([]byte(content))
	f.Close()

	// Should not error because Go json unmarshal overwrites duplicate keys
	_, err = Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEmptyCommand(t *testing.T) {
	content := `{
		"version": 1,
		"services": {
			"api": { "command": [] }
		}
	}`

	f, err := os.CreateTemp("", "taskcapsule-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.Write([]byte(content))
	f.Close()

	_, err = Load(f.Name())
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{
		Version: 1,
	}
	applyDefaults(cfg)

	if cfg.Defaults.GracefulShutdownSeconds != 5 {
		t.Errorf("expected 5, got %d", cfg.Defaults.GracefulShutdownSeconds)
	}
	if cfg.Defaults.HealthTimeoutSeconds != 30 {
		t.Errorf("expected 30, got %d", cfg.Defaults.HealthTimeoutSeconds)
	}
	if cfg.Defaults.BranchPrefix != "task/" {
		t.Errorf("expected task/, got %s", cfg.Defaults.BranchPrefix)
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		cmd     []string
		wantErr bool
	}{
		{[]string{"go", "run", "."}, false},
		{[]string{"pnpm", "dev"}, false},
		{[]string{}, true},
		{[]string{"echo", "hello; rm -rf /"}, true},
		{[]string{"echo", "hello && world"}, true},
		{[]string{"echo", "hello | world"}, true},
	}

	for _, tt := range tests {
		err := validateCommand(tt.cmd)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateCommand(%v) error = %v, wantErr = %v", tt.cmd, err, tt.wantErr)
		}
	}
}
