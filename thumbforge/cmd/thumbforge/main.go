package main

import (
	"fmt"
	"os"

	"github.com/pekomon/go-sandbox/thumbforge/internal/cli"
	"github.com/pekomon/go-sandbox/thumbforge/internal/thumbforge"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	cfg, err := cli.ParseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if _, err := thumbforge.Generate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
