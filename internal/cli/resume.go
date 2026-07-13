package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "resume",
		desc: "Resume a paused capsule and restart its services",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule resume <capsule-name>\n")
				return 2
			}

			opts := app.ResumeOptions{}
			for i := 1; i < len(args); i++ {
				if args[i] == "--setup" {
					opts.RunSetup = true
				}
			}

			result, err := app.Resume(name, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			if result.AlreadyRunning {
				fmt.Printf("Capsule is already running: %s\n", name)
				return 0
			}

			fmt.Printf("Capsule resumed: %s\n\n", name)

			fmt.Println("Services:")
			for _, s := range result.Services {
				printRunning(s.Name, s.PID, s.Port)
			}

			if result.LastNote != "" {
				fmt.Println()
				fmt.Printf("Last note:\n  %s\n", result.LastNote)
			}

			return 0
		},
	})
}
