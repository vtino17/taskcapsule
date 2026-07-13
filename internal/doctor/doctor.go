package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vtino17/taskcapsule/internal/capsule"
	"github.com/vtino17/taskcapsule/internal/state"
)

type CheckResult struct {
	OK      bool
	Message string
}

func Run(stateBase string) ([]CheckResult, error) {
	var results []CheckResult

	// Check git
	if _, err := exec.LookPath("git"); err != nil {
		results = append(results, CheckResult{Message: "Git not found in PATH"})
	} else {
		results = append(results, CheckResult{OK: true, Message: "Git available"})
	}

	// Check state directory
	if err := os.MkdirAll(filepath.Join(stateBase, "capsules"), 0755); err != nil {
		results = append(results, CheckResult{Message: fmt.Sprintf("State directory not writable: %v", err)})
	} else {
		results = append(results, CheckResult{OK: true, Message: "State directory writable"})
	}

	// Check worktree directory
	if err := os.MkdirAll(filepath.Join(stateBase, "worktrees"), 0755); err != nil {
		results = append(results, CheckResult{Message: fmt.Sprintf("Worktree directory not writable: %v", err)})
	} else {
		results = append(results, CheckResult{OK: true, Message: "Worktree directory writable"})
	}

	// Check each capsule
	store := state.NewStore(stateBase)
	capsules, err := store.ListAll()
	if err != nil {
		results = append(results, CheckResult{Message: fmt.Sprintf("Cannot list capsules: %v", err)})
		return results, nil
	}

	for _, c := range capsules {
		checkCapsule(c, &results)
	}

	return results, nil
}

func checkCapsule(c *capsule.State, results *[]CheckResult) {
	// Check worktree
	if _, err := os.Stat(c.WorktreePath); os.IsNotExist(err) {
		*results = append(*results, CheckResult{
			Message: fmt.Sprintf("Capsule %q is missing its worktree: %s", c.Name, c.WorktreePath),
		})
	}

	// Check state file
	stateStore := state.NewStore(filepath.Dir(filepath.Dir(c.WorktreePath)))
	_, err := stateStore.Load(c.RepositoryID, c.Name)
	if err != nil {
		*results = append(*results, CheckResult{
			Message: fmt.Sprintf("Capsule %q has unreadable state: %v", c.Name, err),
		})
	}

	// Check stale PID
	if c.Status == "running" {
		hasLive := false
		for _, svc := range c.Services {
			if svc.PID > 0 {
				proc, err := os.FindProcess(svc.PID)
				if err == nil && proc.Signal(os.Interrupt) == nil {
					hasLive = true
				}
			}
		}
		if !hasLive && len(c.Services) > 0 {
			*results = append(*results, CheckResult{
				Message: fmt.Sprintf("Capsule %q has stale PID (no running services)", c.Name),
			})
		}
	}

	// Check log directory
	logDir := filepath.Join(filepath.Dir(c.WorktreePath), "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Not an error, but worth noting
	}
}
