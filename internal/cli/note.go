package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/vtino17/taskcapsule/internal/app"
)

func init() {
	register(command{
		name: "note",
		desc: "Save a context note for a capsule",
		handler: func(args []string) int {
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "Usage: taskcapsule note <capsule-name> <note-text>\n")
				return 2
			}

			name := args[0]
			text := strings.Join(args[1:], " ")

			if err := app.SaveNote(name, text); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return exitCodeFromError(err)
			}

			fmt.Println("Note saved for " + name + ".")
			return 0
		},
	})
}
