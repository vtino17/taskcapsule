package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "check",
		desc: "Run a validation command in the capsule worktree",
		handler: func(args []string) int {
			if len(args) < 2 || args[1] != "--" {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule check <capsule-name> -- <command> [args...]\n")
				return 2
			}

			name := args[0]
			cmdArgs := args[2:]

			if len(cmdArgs) == 0 {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule check <capsule-name> -- <command> [args...]\n")
				return 2
			}

			result, err := app.RunCheck(name, cmdArgs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			if result.ExitCode == 0 {
				fmt.Printf("Check passed: %s\n\n", name)
			} else {
				fmt.Printf("Check failed: %s\n\n", name)
			}

			fmt.Printf("Command: %s\n", result.Command)
			fmt.Printf("Duration: %.1fs\n", result.Duration)
			fmt.Printf("Exit code: %d\n", result.ExitCode)

			if result.ExitCode != 0 && result.LogPath != "" {
				fmt.Println()
				fmt.Printf("Log:\n  %s\n", result.LogPath)
			}

			return result.ExitCode
		},
	})
}
