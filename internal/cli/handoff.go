package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "handoff",
		desc: "Generate a Markdown handoff report for a capsule",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule handoff <capsule-name>\n")
				return 2
			}

			path, err := app.GenerateHandoff(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Printf("Handoff generated:\n  %s\n", path)
			return 0
		},
	})
}
