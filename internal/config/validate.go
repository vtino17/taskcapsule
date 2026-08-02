package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vtino17/taskcapsule/internal/capsule"
)

var environmentNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func validate(cfg *Config) error {
	if cfg.Version < 1 || cfg.Version > 1 {
		return fmt.Errorf("unsupported schema version: %d (expected 1)", cfg.Version)
	}

	seen := make(map[string]bool)
	for name := range cfg.Services {
		if err := capsule.ValidateName(name); err != nil {
			return fmt.Errorf("invalid service name %q: %v", name, err)
		}
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
		if filepath.IsAbs(svc.WorkingDirectory) {
			return fmt.Errorf("service %q: working directory must be relative to the worktree", name)
		}
		for envName := range svc.Environment {
			if !environmentNamePattern.MatchString(envName) {
				return fmt.Errorf("service %q: invalid environment variable name %q", name, envName)
			}
		}
		for _, envName := range svc.InheritEnvironment {
			if !environmentNamePattern.MatchString(envName) {
				return fmt.Errorf("service %q: invalid inherited environment variable name %q", name, envName)
			}
		}
		if svc.Health != nil {
			switch svc.Health.Type {
			case "none", "process", "tcp", "http":
			default:
				return fmt.Errorf("service %q: unsupported health check type %q", name, svc.Health.Type)
			}
			if svc.Health.TimeoutSeconds < 0 {
				return fmt.Errorf("service %q: health timeout must not be negative", name)
			}
		}
	}

	seenChecks := make(map[string]bool)
	for name := range cfg.Checks {
		if err := capsule.ValidateName(name); err != nil {
			return fmt.Errorf("invalid check name %q: %v", name, err)
		}
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

	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve worktree root %q: %w", worktreeRoot, err)
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return fmt.Errorf("resolve working directory %q: %w", dir, err)
	}

	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return fmt.Errorf("working directory %q is outside worktree", dir)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
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
