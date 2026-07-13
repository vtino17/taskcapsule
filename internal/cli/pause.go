package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "pause",
		desc: "Stop all services in a capsule and release resources",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule pause <capsule-name>\n")
				return 2
			}

			result, err := app.Pause(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			if result.AlreadyPaused {
				fmt.Printf("Capsule is already paused: %s\n", name)
				return 0
			}

			fmt.Printf("Capsule paused: %s\n\n", name)

			fmt.Println("Stopped:")
			for _, s := range result.Services {
				printStopped(s.Name, s.PID)
			}

			fmt.Println()
			fmt.Println("Resources released.")
			return 0
		},
	})
}
