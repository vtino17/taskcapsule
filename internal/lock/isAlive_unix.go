//go:build !windows

package lock

import (
	"os"
	"syscall"
)

func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 is a no-op that just checks process existence
	return proc.Signal(syscall.Signal(0)) == nil
}
