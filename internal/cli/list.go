package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "list",
		desc: "List all capsules in the current repository",
		handler: func(args []string) int {
			showAll := false
			for _, a := range args {
				if a == "--all" || a == "-a" {
					showAll = true
				}
			}

			capsules, err := app.List(showAll)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return 1
			}

			if len(capsules) == 0 {
				fmt.Println("No capsules found.")
				return 0
			}

			fmt.Printf("%-20s %-12s %-28s %s\n", "NAME", "STATUS", "BRANCH", "UPDATED")
			for _, c := range capsules {
				fmt.Printf("%-20s %-12s %-28s %s\n", c.Name, c.Status, c.Branch, c.Updated)
			}
			return 0
		},
	})
}
