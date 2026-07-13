package cli

import (
	"fmt"
	"os"
	"strings"
)

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "capsule not found"):
		return 3
	case strings.Contains(msg, "not a git repository"):
		return 5
	case strings.Contains(msg, "already exists"):
		return 4
	case strings.Contains(msg, "cannot delete running"):
		return 4
	case strings.Contains(msg, "uncommitted changes"):
		return 4
	case strings.Contains(msg, "already paused"):
		return 0
	case strings.Contains(msg, "already running"):
		return 0
	default:
		return 1
	}
}

var noColor = os.Getenv("NO_COLOR") != ""

const (
	markSuccess = "\u2713"
	markWarning = "!"
	markFailure = "\u00d7"
	markRunning = "\u25cf"
	markStopped = "\u25cb"
)

func printSuccess(format string, a ...any) {
	fmt.Printf("%s %s\n", markSuccess, fmt.Sprintf(format, a...))
}

func printWarning(format string, a ...any) {
	fmt.Printf("%s %s\n", markWarning, fmt.Sprintf(format, a...))
}

func printFailure(format string, a ...any) {
	fmt.Printf("%s %s\n", markFailure, fmt.Sprintf(format, a...))
}

func printRunning(name string, pid, port int) {
	if port > 0 {
		fmt.Printf("  %s  %-15s pid=%-6d port=%d\n", markRunning, name, pid, port)
	} else {
		fmt.Printf("  %s  %-15s pid=%d\n", markRunning, name, pid)
	}
}

func printStopped(name string, pid int) {
	fmt.Printf("  %s  %-15s pid=%d\n", markStopped, name, pid)
}
