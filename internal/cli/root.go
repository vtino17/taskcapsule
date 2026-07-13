package cli

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/vtino17/taskcapsule/internal/version"
)

type command struct {
	name    string
	desc    string
	handler func(args []string) int
}

var commands []command

func register(cmd command) {
	commands = append(commands, cmd)
}

func Run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 2
	}

	name := args[0]
	for _, cmd := range commands {
		if cmd.name == name {
			return cmd.handler(args[1:])
		}
	}

	switch name {
	case "--help", "-h":
		printUsage()
		return 0
	case "--version", "-v":
		fmt.Println(version.String())
		return 0
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", name)
	fmt.Fprintf(os.Stderr, "Run 'taskcapsule --help' for usage.\n")
	return 2
}

func printUsage() {
	fmt.Println("TaskCapsule - Pause and resume coding tasks without losing your place.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  taskcapsule <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")

	max := 0
	for _, cmd := range commands {
		if len(cmd.name) > max {
			max = len(cmd.name)
		}
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].name < commands[j].name
	})
	for _, cmd := range commands {
		padding := strings.Repeat(" ", max-len(cmd.name)+2)
		fmt.Printf("  %s%s%s\n", cmd.name, padding, cmd.desc)
	}

	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help, -h     Show this help")
	fmt.Println("  --version, -v  Show version")
	fmt.Println()
	fmt.Println("Use 'taskcapsule <command> --help' for command-specific help.")
}

func requireArgs(args []string, min int, usage string) error {
	if len(args) < min {
		return errors.New(usage)
	}
	return nil
}

func requireCapsuleName(args []string) (string, error) {
	if err := requireArgs(args, 1, "Usage: taskcapsule <command> <capsule-name>"); err != nil {
		return "", err
	}
	return args[0], nil
}
