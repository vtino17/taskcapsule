package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	if cfg.Version == 0 {
		cfg.Version = 1
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Defaults.GracefulShutdownSeconds <= 0 {
		cfg.Defaults.GracefulShutdownSeconds = 5
	}
	if cfg.Defaults.HealthTimeoutSeconds <= 0 {
		cfg.Defaults.HealthTimeoutSeconds = 30
	}
	if cfg.Defaults.BranchPrefix == "" {
		cfg.Defaults.BranchPrefix = "task/"
	}
}
