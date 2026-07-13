package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "where",
		desc: "Show a summary to continue working on a capsule",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule where <capsule-name>\n")
				return 2
			}

			w, err := app.Where(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Printf("You were working on: %s\n\n", w.Name)

			if w.LastNote != "" {
				fmt.Printf("Last note:\n  %s\n\n", w.LastNote)
			}

			if len(w.ModifiedFiles) > 0 {
				fmt.Println("Last modified files:")
				for _, f := range w.ModifiedFiles {
					fmt.Printf("  %s\n", f)
				}
				fmt.Println()
			}

			if w.LastCheck != nil {
				fmt.Println("Last check:")
				fmt.Printf("  %s\n", w.LastCheck.Command)
				if w.LastCheck.ExitCode == 0 {
					fmt.Println("  Result: passed")
				} else {
					fmt.Println("  Result: failed")
				}
				fmt.Println()
			}

			fmt.Println("Suggested next action:")
			if w.Status == "paused" {
				fmt.Printf("  Resume the capsule: taskcapsule resume %s\n", name)
			} else {
				fmt.Printf("  Go to the worktree and continue working.\n")
			}

			return 0
		},
	})
}
