//go:build windows

package process

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

var _ = syscall.SIGTERM

type sysProcAttr = syscall.SysProcAttr

func setProcessGroup(cmd *exec.Cmd) {
	// Windows: Job Objects would be used here. For MVP, no special setup.
}

func getProcessGroup(cmd *exec.Cmd) int {
	if cmd.Process != nil {
		return cmd.Process.Pid
	}
	return 0
}

func stopProcessGroup(pgid int, graceSeconds int) {
	// Windows MVP: signal the main process
	proc, err := os.FindProcess(pgid)
	if err != nil {
		return
	}
	proc.Signal(os.Interrupt)
	time.Sleep(time.Duration(graceSeconds) * time.Second)
	proc.Kill()
}

func stopProcess(pid int, graceSeconds int) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	proc.Signal(os.Interrupt)
	time.Sleep(time.Duration(graceSeconds) * time.Second)
	proc.Kill()
}

func isAlive(proc *os.Process) bool {
	// Windows: FindProcess always succeeds, so we can't reliably check
	// This is a known limitation for MVP
	return true
}
