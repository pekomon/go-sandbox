package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("memesweeper", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var showVersion bool
	fs.BoolVar(&showVersion, "version", false, "print the MemeSweeper version")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if showVersion {
		fmt.Fprintln(os.Stdout, version)
		return 0
	}

	fmt.Fprintln(os.Stderr, "MemeSweeper gameplay is not implemented yet. Follow upcoming issues for progress.")
	return 1
}
