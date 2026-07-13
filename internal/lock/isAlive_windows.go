//go:build windows

package lock

import "os"

func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, os.FindProcess returns error if process doesn't exist.
	// If it succeeds, send a signal to verify.
	err = proc.Signal(os.Interrupt)
	return err == nil
}
