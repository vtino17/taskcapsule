package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "status",
		desc: "Show detailed capsule status",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule status <capsule-name>\n")
				return 2
			}

			s, err := app.Status(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Printf("Capsule: %s\n", s.Name)
			fmt.Printf("Status: %s\n\n", s.Status)

			fmt.Printf("Repository:\n  %s\n\n", s.RepositoryRoot)

			fmt.Printf("Worktree:\n  %s\n\n", s.WorktreePath)

			fmt.Println("Git:")
			fmt.Printf("  Branch: %s\n", s.Branch)
			fmt.Printf("  Base: %s\n", s.BaseBranch)
			fmt.Printf("  Dirty: %s\n", s.Dirty)
			fmt.Printf("  Changed files: %d\n\n", s.ChangedFiles)

			if len(s.Services) > 0 {
				fmt.Println("Services:")
				for _, svc := range s.Services {
					if svc.Running {
						printRunning(svc.Name, svc.PID, svc.Port)
					} else {
						printStopped(svc.Name, svc.PID)
					}
				}
				fmt.Println()
			}

			if s.LastCheck != nil {
				fmt.Println("Last check:")
				fmt.Printf("  %s\n", s.LastCheck.Command)
				if s.LastCheck.ExitCode == 0 {
					fmt.Println("  Passed")
				} else {
					fmt.Println("  Failed")
				}
				fmt.Printf("  Exit code: %d\n", s.LastCheck.ExitCode)
				fmt.Printf("  Finished: %s\n\n", s.LastCheck.FinishedAt)
			}

			if s.LastNote != "" {
				fmt.Printf("Last note:\n  %s\n", s.LastNote)
			}

			return 0
		},
	})
}
