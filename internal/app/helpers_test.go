package app

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/ports"
)

func TestBuildServiceEnvAllocatesEachPlaceholderOnce(t *testing.T) {
	cfg := config.ServiceConfig{Environment: map[string]string{
		"PORT": "${PORT:api}",
		"URL":  "http://127.0.0.1:${PORT:api}/health",
	}}

	env, err := buildServiceEnv(cfg, ports.NewAllocator())
	if err != nil {
		t.Fatal(err)
	}
	if len(env.Ports) != 1 || env.Ports["api"] <= 0 {
		t.Fatalf("expected one allocated api port, got %#v", env.Ports)
	}
}

func TestSetupEnvUsesExplicitInheritance(t *testing.T) {
	t.Setenv("TASKCAPSULE_TEST_SECRET", "must-not-leak")
	t.Setenv("TASKCAPSULE_TEST_ALLOWED", "allowed")

	cfg := config.ServiceConfig{
		Environment:        map[string]string{"APP_MODE": "test", "PORT": "${PORT:api}"},
		InheritEnvironment: []string{"TASKCAPSULE_TEST_ALLOWED"},
	}
	cmd := exec.Command("ignored")
	setupEnv(cmd, cfg, &serviceEnv{Ports: map[string]int{"api": 43210}})
	joined := strings.Join(cmd.Env, "\n")

	if strings.Contains(joined, "TASKCAPSULE_TEST_SECRET=") {
		t.Fatal("unlisted environment variable leaked to child process")
	}
	for _, expected := range []string{
		"TASKCAPSULE_TEST_ALLOWED=allowed",
		"APP_MODE=test",
		"PORT=43210",
		"PORT_API=43210",
	} {
		if !strings.Contains(joined, expected) {
			t.Errorf("expected %q in child environment: %v", expected, cmd.Env)
		}
	}
}
