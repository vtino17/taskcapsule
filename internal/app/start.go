package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/config"
	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/health"
	"github.com/vtino17/taskcapsule/internal/ports"
	"github.com/vtino17/taskcapsule/internal/process"
	"github.com/vtino17/taskcapsule/internal/state"
)

type startedService struct {
	Name    string
	PID     int
	PGID    int
	Port    int
	LogPath string
}

func Start(name string, opts StartOptions) (*StartResult, error) {
	root, err := findGitRoot()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(filepath.Join(root, ".taskcapsule.json"))
	if err != nil {
		return nil, err
	}

	if err := capsule.ValidateName(name); err != nil {
		return nil, err
	}

	slug := capsule.Slugify(name)

	repoID, err := git.RepoID(root)
	if err != nil {
		return nil, err
	}

	repoName, err := git.RepoName(root)
	if err != nil {
		return nil, err
	}

	stateBase, err := getStateDir()
	if err != nil {
		return nil, err
	}

	cs := state.NewStore(stateBase)

	existing, _ := cs.Load(repoID, slug)
	if existing != nil {
		return nil, fmt.Errorf("capsule %q already exists in this repository", name)
	}

	// Acquire lock before mutations
	cl, err := acquireCapsuleLock(repoID, slug, "start")
	if err != nil {
		return nil, err
	}
	defer cl.Release()

	baseBranch := opts.BaseBranch
	if baseBranch == "" {
		baseBranch = cfg.Defaults.BaseBranch
		if baseBranch == "" {
			baseBranch, _ = git.DefaultBranch(root)
		}
		if baseBranch == "" {
			baseBranch = "main"
		}
	}

	capsuleBranch := opts.Branch
	if capsuleBranch == "" {
		prefix := cfg.Defaults.BranchPrefix
		if prefix == "" {
			prefix = "task/"
		}
		capsuleBranch = prefix + name
	}

	worktreesDir := filepath.Join(stateBase, "worktrees")
	worktreePath := filepath.Join(worktreesDir, repoName, slug)

	inUse, _ := git.BranchInUse(capsuleBranch, root)
	if inUse {
		return nil, fmt.Errorf("branch %q is already in use by another worktree", capsuleBranch)
	}

	now := time.Now().UTC()
	newState := &capsule.State{
		SchemaVersion:  1,
		Name:           slug,
		Status:         "preparing",
		RepositoryRoot: root,
		RepositoryID:   repoID,
		WorktreePath:   worktreePath,
		Branch:         capsuleBranch,
		BaseBranch:     baseBranch,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := cs.Save(repoID, slug, newState); err != nil {
		return nil, fmt.Errorf("failed to save initial state: %v", err)
	}

	// Create worktree
	if err := git.CreateWorktree(root, worktreePath, capsuleBranch, baseBranch); err != nil {
		newState.Status = "error"
		newState.LastError = strPtr(fmt.Sprintf("failed to create worktree: %v", err))
		newState.UpdatedAt = time.Now().UTC()
		cs.Save(repoID, slug, newState)
		return nil, fmt.Errorf("failed to create worktree: %v\nPath: %s", err, worktreePath)
	}

	// Run setup commands
	for _, setup := range cfg.Setup {
		cmd := execInDir(setup.Command, worktreePath)
		if output, err := cmd.CombinedOutput(); err != nil {
			newState.Status = "error"
			newState.LastError = strPtr(fmt.Sprintf("setup failed: %v", err))
			newState.UpdatedAt = time.Now().UTC()
			cs.Save(repoID, slug, newState)
			return nil, fmt.Errorf("setup command failed: %s\n%s", strings.Join(setup.Command, " "), string(output))
		}
	}

	if opts.NoServices {
		newState.Status = "running"
		newState.UpdatedAt = time.Now().UTC()
		cs.Save(repoID, slug, newState)
		return &StartResult{
			Name:         slug,
			Branch:       capsuleBranch,
			WorktreePath: worktreePath,
			Status:       "running",
		}, nil
	}

	// Sequential service startup with health checks and rollback
	started, serviceInfos, startErr := startServices(cfg, stateBase, repoID, slug, worktreePath, newState, cs)
	if startErr != nil {
		rollbackServices(started, newState, cs)
		setErrorState(newState, startErr, cs, repoID, slug)
		return nil, startErr
	}

	newState.Status = "running"
	newState.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, slug, newState)

	return &StartResult{
		Name:         slug,
		Branch:       capsuleBranch,
		WorktreePath: worktreePath,
		Status:       "running",
		Services:     serviceInfos,
	}, nil
}

func startServices(cfg *config.Config, stateBase, repoID, slug, worktreePath string, newState *capsule.State, cs *state.Store) ([]startedService, []ServiceInfo, error) {
	allocator := ports.NewAllocator()
	var started []startedService

	for svcName, svcCfg := range cfg.Services {
		env := buildServiceEnv(svcCfg, allocator)

		logDir := filepath.Join(stateBase, "capsules", repoID, slug, "logs")
		os.MkdirAll(logDir, 0755)
		logFile := filepath.Join(logDir, svcName+".log")

		f, err := os.Create(logFile)
		if err != nil {
			return started, nil, fmt.Errorf("cannot create log file for %s: %v", svcName, err)
		}

		cmd := execInDir(svcCfg.Command, worktreePath)
		cmd.Stdout = f
		cmd.Stderr = f
		process.SetProcessGroup(cmd)

		if err := cmd.Start(); err != nil {
			f.Close()
			return started, nil, fmt.Errorf("cannot start service %s: %v\nLog: %s", svcName, err, logFile)
		}

		pgid := process.GetProcessGroup(cmd)
		pid := cmd.Process.Pid
		port := env.Ports[svcName]

		svc := startedService{
			Name:    svcName,
			PID:     pid,
			PGID:    pgid,
			Port:    port,
			LogPath: logFile,
		}
		started = append(started, svc)

		svcState := capsule.ServiceState{
			Status:       "starting",
			Command:      svcCfg.Command,
			PID:          pid,
			ProcessGroup: pgid,
			Port:         port,
			LogPath:      logFile,
		}
		if newState.Services == nil {
			newState.Services = make(map[string]capsule.ServiceState)
		}
		newState.Services[svcName] = svcState
		newState.UpdatedAt = time.Now().UTC()
		cs.Save(repoID, slug, newState)

		// Health check
		if svcCfg.Health != nil {
			hcCfg := health.Config{
				Type:           svcCfg.Health.Type,
				URL:            resolvePortVar(svcCfg.Health.URL, env.Ports),
				Host:           svcCfg.Health.Host,
				Port:           resolvePortVar(svcCfg.Health.Port, env.Ports),
				ExpectedStatus: svcCfg.Health.ExpectedStatus,
				TimeoutSeconds: svcCfg.Health.TimeoutSeconds,
				PID:            pid,
				ProcessGroup:   pgid,
			}

			if hcCfg.TimeoutSeconds <= 0 {
				hcCfg.TimeoutSeconds = cfg.Defaults.HealthTimeoutSeconds
			}

			result := health.Check(hcCfg)
			if !result.OK {
				return started, nil, fmt.Errorf(
					"service %s failed health check\n\nHealth type: %s\nReason: %s\nLog: %s",
					svcName, hcCfg.Type, result.Error, logFile)
			}
		}

		svcState.Status = "running"
		newState.Services[svcName] = svcState
		newState.UpdatedAt = time.Now().UTC()
		cs.Save(repoID, slug, newState)
	}

	serviceInfos := make([]ServiceInfo, 0, len(started))
	for _, s := range started {
		serviceInfos = append(serviceInfos, ServiceInfo{
			Name:    s.Name,
			PID:     s.PID,
			Port:    s.Port,
			Running: true,
		})
	}

	return started, serviceInfos, nil
}

func rollbackServices(started []startedService, state *capsule.State, cs *state.Store) {
	repoID := state.RepositoryID

	for i := len(started) - 1; i >= 0; i-- {
		s := started[i]
		if s.PGID > 0 {
			process.StopProcessGroup(s.PGID, 5)
		} else if s.PID > 0 {
			process.StopProcess(s.PID, 5)
		}

		if svc, ok := state.Services[s.Name]; ok {
			svc.Status = "stopped"
			svc.PID = 0
			svc.ProcessGroup = 0
			state.Services[s.Name] = svc
		}
	}

	state.Status = "error"
	state.UpdatedAt = time.Now().UTC()
	cs.Save(repoID, state.Name, state)
}

func resolvePortVar(s string, ports map[string]int) string {
	if s == "" {
		return s
	}
	result := s
	for name, port := range ports {
		placeholder := fmt.Sprintf("${PORT:%s}", name)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%d", port))
	}
	return result
}
