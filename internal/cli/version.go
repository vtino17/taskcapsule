package cli

import (
	"fmt"

	"github.com/vtino17/taskcapsule/internal/version"
)

func init() {
	register(command{
		name: "version",
		desc: "Show version information",
		handler: func(args []string) int {
			fmt.Println(version.String())
			return 0
		},
	})
}
