package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "logs",
		desc: "Show service logs for a capsule",
		handler: func(args []string) int {
			if len(args) < 1 {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule logs <capsule-name> [service-name] [--lines N] [--follow]\n")
				return 2
			}

			name := args[0]
			opts := app.LogOptions{Lines: 50}

			for i := 1; i < len(args); i++ {
				switch args[i] {
				case "--lines":
					i++
					if i < len(args) {
						if n, err := strconv.Atoi(args[i]); err == nil {
							opts.Lines = n
						}
					}
				case "--follow", "-f":
					opts.Follow = true
				default:
					opts.ServiceName = args[i]
				}
			}

			logs, err := app.ShowLogs(name, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			os.Stdout.Write(logs)
			return 0
		},
	})
}
