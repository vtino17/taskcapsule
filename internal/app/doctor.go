package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/vtino17/taskcapsule/internal/git"
	"github.com/vtino17/taskcapsule/internal/state"
)

func Doctor() ([]DoctorResult, error) {
	var results []DoctorResult

	// Check git available
	if _, err := exec.LookPath("git"); err != nil {
		results = append(results, DoctorResult{OK: false, Message: "Git not found in PATH"})
	} else {
		results = append(results, DoctorResult{OK: true, Message: "Git available"})
	}

	// Check we're in a git repo
	root, err := findGitRoot()
	if err != nil {
		results = append(results, DoctorResult{OK: false, Message: "Not in a git repository"})
		return results, nil
	}

	// Check config
	cfgPath := filepath.Join(root, ".taskcapsule.json")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		results = append(results, DoctorResult{OK: false, Message: "Configuration not found: run 'taskcapsule init'"})
	} else {
		results = append(results, DoctorResult{OK: true, Message: "Configuration valid"})
	}

	// Check state directory
	stateBase, err := getStateDir()
	if err != nil {
		results = append(results, DoctorResult{OK: false, Message: "Cannot determine state directory"})
	} else {
		if err := os.MkdirAll(filepath.Join(stateBase, "capsules"), 0755); err != nil {
			results = append(results, DoctorResult{OK: false, Message: "State directory not writable"})
		} else {
			results = append(results, DoctorResult{OK: true, Message: "State directory writable"})
		}

		if err := os.MkdirAll(filepath.Join(stateBase, "worktrees"), 0755); err != nil {
			results = append(results, DoctorResult{OK: false, Message: "Worktree directory not writable"})
		} else {
			results = append(results, DoctorResult{OK: true, Message: "Worktree directory writable"})
		}
	}

	// Check all capsules
	cs := state.NewStore(stateBase)
	capsules, err := cs.ListAll()
	if err == nil {
		for _, caps := range capsules {
			// Check worktree exists
			if _, err := os.Stat(caps.WorktreePath); os.IsNotExist(err) {
				results = append(results, DoctorResult{
					OK:      false,
					Message: "Capsule " + caps.Name + " is missing its worktree",
				})
				continue
			}

			// Check branch still exists
			branchExists, _ := git.BranchExists(caps.Branch, root)
			if !branchExists {
				results = append(results, DoctorResult{
					OK:      false,
					Message: "Capsule " + caps.Name + " branch '" + caps.Branch + "' no longer exists",
				})
			}

			// Check for stale PIDs
			if caps.Status == "running" {
				hasLiveProcess := false
				for _, svc := range caps.Services {
					if svc.PID > 0 && isProcessRunning(svc.PID) {
						hasLiveProcess = true
					}
				}
				if !hasLiveProcess && len(caps.Services) > 0 {
					results = append(results, DoctorResult{
						OK:      false,
						Message: "Capsule " + caps.Name + " has stale PID (services not running)",
					})
				}
			}
		}
	}

	return results, nil
}

func isProcessRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks process existence on Unix.
	// On Windows, FindProcess always succeeds, so this returns true.
	return proc.Signal(syscall.Signal(0x0)) == nil
}
