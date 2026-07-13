package main

import (
	"fmt"
	"os"

	"github.com/vtino17/taskcapsule/internal/cli"
)

func main() {
	code := cli.Run(os.Args[1:])
	os.Exit(code)
}

func init() {
	// Recovered panics become exit code 10 (internal error)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			os.Exit(10)
		}
	}()
}
