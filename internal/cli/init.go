package cli

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/config"
)

func init() {
	register(command{
		name: "init",
		desc: "Create initial configuration file",
		handler: func(args []string) int {
			force := false
			for _, a := range args {
				if a == "--force" || a == "-f" {
					force = true
				}
			}

			wd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot get working directory: %v\n", err)
				return 1
			}

			path := wd + "/.taskcapsule.json"

			if !force {
				if _, err := os.Stat(path); err == nil {
					fmt.Fprintf(os.Stderr, "Configuration already exists: .taskcapsule.json\n")
					fmt.Fprintf(os.Stderr, "Use --force to replace it.\n")
					return 2
				}
			}

			data := config.DefaultTemplate()
			if err := os.WriteFile(path, []byte(data), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write configuration: %v\n", err)
				return 1
			}

			fmt.Println("Created .taskcapsule.json")
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("1. Review the generated configuration")
			fmt.Println("2. Run: taskcapsule start my-task")
			return 0
		},
	})
}
