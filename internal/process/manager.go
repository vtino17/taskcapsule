package process

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ProcessInfo struct {
	PID          int
	ProcessGroup int
}

// SetProcessGroup configures a command to run in a new process group.
// Unix: uses Setpgid. Windows: no-op (experimental).
func SetProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &sysProcAttr{}
	}
	setProcessGroup(cmd)
}

// GetProcessGroup returns the process group ID of a command.
func GetProcessGroup(cmd *exec.Cmd) int {
	return getProcessGroup(cmd)
}

// StopProcessGroup sends SIGTERM then SIGKILL to a process group.
func StopProcessGroup(pgid int, graceSeconds int) {
	if pgid <= 0 {
		return
	}
	stopProcessGroup(pgid, graceSeconds)
}

// StopProcess sends SIGTERM then SIGKILL to a single process.
func StopProcess(pid int, graceSeconds int) {
	if pid <= 0 {
		return
	}
	stopProcess(pid, graceSeconds)
}

// IsAlive checks if a process is still running.
func IsAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return isAlive(proc)
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.1fs", d.Seconds())
}
