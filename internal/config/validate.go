package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validate(cfg *Config) error {
	if cfg.Version < 1 || cfg.Version > 1 {
		return fmt.Errorf("unsupported schema version: %d (expected 1)", cfg.Version)
	}

	seen := make(map[string]bool)
	for name := range cfg.Services {
		if seen[name] {
			return fmt.Errorf("duplicate service name: %s", name)
		}
		seen[name] = true

		svc := cfg.Services[name]
		if len(svc.Command) == 0 {
			return fmt.Errorf("service %q has empty command", name)
		}
		if err := validateCommand(svc.Command); err != nil {
			return fmt.Errorf("service %q: %v", name, err)
		}
	}

	seenChecks := make(map[string]bool)
	for name := range cfg.Checks {
		if seenChecks[name] {
			return fmt.Errorf("duplicate check name: %s", name)
		}
		seenChecks[name] = true

		chk := cfg.Checks[name]
		if len(chk.Command) == 0 {
			return fmt.Errorf("check %q has empty command", name)
		}
	}

	for i, setup := range cfg.Setup {
		if len(setup.Command) == 0 {
			return fmt.Errorf("setup command %d has empty command", i)
		}
	}

	return nil
}

func validateCommand(cmd []string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command must not be empty")
	}
	for _, part := range cmd {
		if strings.Contains(part, "..") || strings.Contains(part, "|") || strings.Contains(part, ";") || strings.Contains(part, "&&") {
			return fmt.Errorf("command part %q contains shell metacharacters", part)
		}
	}
	return nil
}

func ValidateWorkDir(dir, worktreeRoot string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	root, err := filepath.Abs(worktreeRoot)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return fmt.Errorf("working directory %q is outside worktree", dir)
	}
	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("working directory %q is outside worktree root %q", dir, root)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("working directory %q is not a directory", dir)
	}

	return nil
}
