package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "doctor",
		desc: "Check TaskCapsule installation and capsule states",
		handler: func(args []string) int {
			results, err := app.Doctor()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return 1
			}

			fmt.Println("TaskCapsule Doctor")
			fmt.Println()

			issues := 0
			for _, r := range results {
				mark := markSuccess
				if !r.OK {
					mark = markWarning
					issues++
				}
				fmt.Printf("%s %s\n", mark, r.Message)
			}

			fmt.Println()
			if issues > 0 {
				fmt.Printf("%d issue(s) detected.\n", issues)
				return 1
			}
			fmt.Println("All checks passed.")
			return 0
		},
	})
}
