package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "delete",
		desc: "Delete a capsule and its worktree",
		handler: func(args []string) int {
			name, err := requireCapsuleName(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule delete <capsule-name> [--force]\n")
				return 2
			}

			force := false
			for _, a := range args[1:] {
				if a == "--force" || a == "-f" {
					force = true
				}
			}

			if err := app.DeleteCapsule(name, force); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Printf("Capsule deleted: %s\n", name)
			return 0
		},
	})
}
