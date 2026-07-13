package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/ports"
	"github.com/vtino17/taskcapsule/internal/state"
)

func execInDir(command []string, dir string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	return cmd
}

type serviceEnv struct {
	Ports map[string]int
}

func buildServiceEnv(svcCfg config.ServiceConfig, allocator *ports.Allocator) *serviceEnv {
	env := &serviceEnv{Ports: make(map[string]int)}

	// Allocate ports for this service - scan all env vars for ${PORT:name}
	for key, val := range svcCfg.Environment {
		_ = key
		if strings.HasPrefix(val, "${PORT:") && strings.HasSuffix(val, "}") {
			name := val[7 : len(val)-1]
			if _, ok := env.Ports[name]; !ok {
				port, _ := allocator.Allocate()
				env.Ports[name] = port
			}
		}
	}

	if portVar, ok := svcCfg.Environment["PORT"]; ok {
		if strings.HasPrefix(portVar, "${PORT:") && strings.HasSuffix(portVar, "}") {
			name := portVar[7 : len(portVar)-1]
			port, _ := allocator.Allocate()
			env.Ports[name] = port
		}
	}

	return env
}

func strPtr(s string) *string {
	return &s
}

func setErrorState(state *capsule.State, err error, cs *state.Store, repoID, name string) {
	state.Status = "error"
	state.LastError = strPtr(err.Error())
	state.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, name, state)
}

func setupEnv(cmd *exec.Cmd, svcEnv *serviceEnv) {
	cmd.Env = os.Environ()

	if svcEnv == nil {
		return
	}

	for svcName, port := range svcEnv.Ports {
		key := fmt.Sprintf("PORT_%s", strings.ToUpper(svcName))
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", key, port))
	}
}

func exitCodeFromError(err error) int {
	// Simple heuristics based on error message
	msg := err.Error()
	if strings.Contains(msg, "capsule not found") {
		return ExitNotFound
	}
	if strings.Contains(msg, "not a git repository") || strings.Contains(msg, "not found") {
		return ExitDependency
	}
	if strings.Contains(msg, "uncommitted changes") || strings.Contains(msg, "already exists") {
		return ExitUnsafe
	}
	if strings.Contains(msg, "already paused") || strings.Contains(msg, "already running") {
		return ExitSuccess
	}
	return ExitFailure
}
