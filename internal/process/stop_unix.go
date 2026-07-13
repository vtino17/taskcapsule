//go:build !windows

package process

import (
	"os"
	"syscall"
	"time"
)

// Additional Unix-specific process management
func waitForExit(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return true
		}
		if proc.Signal(syscall.Signal(0)) != nil {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}
