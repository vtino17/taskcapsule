package app

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/ports"
	"github.com/vtino17/taskcapsule/internal/state"
)

var portPlaceholderPattern = regexp.MustCompile(`\$\{PORT:([A-Za-z_][A-Za-z0-9_]*)\}`)

func execInDir(command []string, dir string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	return cmd
}

type serviceEnv struct {
	Ports map[string]int
}

func buildServiceEnv(svcCfg config.ServiceConfig, allocator *ports.Allocator) (*serviceEnv, error) {
	env := &serviceEnv{Ports: make(map[string]int)}

	for _, value := range svcCfg.Environment {
		for _, match := range portPlaceholderPattern.FindAllStringSubmatch(value, -1) {
			name := match[1]
			if _, ok := env.Ports[name]; ok {
				continue
			}
			port, err := allocator.Allocate()
			if err != nil {
				return nil, err
			}
			env.Ports[name] = port
		}
	}

	return env, nil
}

func strPtr(s string) *string {
	return &s
}

func setErrorState(state *capsule.State, err error, cs *state.Store, repoID, name string) error {
	state.Status = "error"
	state.LastError = strPtr(err.Error())
	state.UpdatedAt = time.Now().UTC()
	return cs.Save(repoID, name, state)
}

func setupEnv(cmd *exec.Cmd, svcCfg config.ServiceConfig, svcEnv *serviceEnv) {
	env := make(map[string]string)
	for _, name := range []string{"PATH", "HOME", "USER", "LOGNAME", "TMPDIR", "TEMP", "TMP", "SystemRoot", "COMSPEC", "PATHEXT"} {
		if value, ok := os.LookupEnv(name); ok {
			env[name] = value
		}
	}
	for _, name := range svcCfg.InheritEnvironment {
		if value, ok := os.LookupEnv(name); ok {
			env[name] = value
		}
	}
	for name, value := range svcCfg.Environment {
		env[name] = resolvePortVar(value, svcEnv.Ports)
	}
	for svcName, port := range svcEnv.Ports {
		key := fmt.Sprintf("PORT_%s", strings.ToUpper(svcName))
		env[key] = fmt.Sprintf("%d", port)
	}

	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	cmd.Env = make([]string, 0, len(keys))
	for _, key := range keys {
		cmd.Env = append(cmd.Env, key+"="+env[key])
	}
}

func exitCodeFromError(err error) int {
	if err == nil {
		return ExitSuccess
	}
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
