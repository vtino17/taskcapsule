package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "start",
		desc: "Create and start a new capsule",
		handler: func(args []string) int {
			if err := requireArgs(args, 1, "Usage: taskcapsule start <capsule-name> [--base <branch>] [--branch <name>] [--no-services]"); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return 2
			}

			name := args[0]
			opts := app.StartOptions{}

			for i := 1; i < len(args); i++ {
				switch args[i] {
				case "--base":
					i++
					if i < len(args) {
						opts.BaseBranch = args[i]
					}
				case "--branch":
					i++
					if i < len(args) {
						opts.Branch = args[i]
					}
				case "--no-services":
					opts.NoServices = true
				}
			}

			result, err := app.Start(name, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Printf("Capsule started: %s\n\n", result.Name)
			fmt.Printf("Branch:    %s\n", result.Branch)
			fmt.Printf("Worktree:  %s\n", result.WorktreePath)
			fmt.Printf("Status:    %s\n\n", result.Status)

			if len(result.Services) > 0 {
				fmt.Println("Services:")
				for _, s := range result.Services {
					printRunning(s.Name, s.PID, s.Port)
				}
				fmt.Println()
			}

			fmt.Printf("Use:\n  taskcapsule pause %s\n", result.Name)
			return 0
		},
	})
}
