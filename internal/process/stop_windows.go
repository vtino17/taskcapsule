//go:build windows

package process

import (
	"time"
)

func waitForExit(pid int, timeout time.Duration) bool {
	// Windows MVP: minimal implementation
	time.Sleep(timeout)
	return true
}
