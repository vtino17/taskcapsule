package config

import (
	"os"
	"path/filepath"
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

func TestLoadRejectsUnsafeMapKeys(t *testing.T) {
	tests := []string{
		`{"version":1,"services":{"../../escape":{"command":["go","run","."]}}}`,
		`{"version":1,"checks":{"../escape":{"command":["go","test","./..."]}}}`,
	}

	for _, content := range tests {
		f, err := os.CreateTemp("", "taskcapsule-*.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			f.Close()
			os.Remove(f.Name())
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			os.Remove(f.Name())
			t.Fatal(err)
		}

		_, loadErr := Load(f.Name())
		os.Remove(f.Name())
		if loadErr == nil {
			t.Fatalf("Load accepted unsafe config: %s", content)
		}
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

func TestValidateWorkDirRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := ValidateWorkDir(link, root); err == nil {
		t.Fatal("working-directory symlink escape was accepted")
	}
}
