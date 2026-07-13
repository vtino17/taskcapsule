//go:build !windows

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
	cmd.SysProcAttr.Setpgid = true
}

func getProcessGroup(cmd *exec.Cmd) int {
	if cmd.Process != nil {
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			return pgid
		}
	}
	return cmd.ProcessState.Pid()
}

func stopProcessGroup(pgid int, graceSeconds int) {
	syscall.Kill(-pgid, syscall.SIGTERM)
	time.Sleep(time.Duration(graceSeconds) * time.Second)
	syscall.Kill(-pgid, syscall.SIGKILL)
}

func stopProcess(pid int, graceSeconds int) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	proc.Signal(syscall.SIGTERM)
	time.Sleep(time.Duration(graceSeconds) * time.Second)
	proc.Signal(syscall.SIGKILL)
}

func isAlive(proc *os.Process) bool {
	return proc.Signal(syscall.Signal(0)) == nil
}
