package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pekomon/go-sandbox/filesort/internal/sorter"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	var dryRun bool
	fs := flag.NewFlagSet("filesort", flag.ContinueOnError)
	fs.BoolVar(&dryRun, "dry-run", false, "plan only; do not modify the filesystem")
	// silence default usage on parse error
	fs.SetOutput(new(nopWriter))
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "invalid flags")
		return 2
	}
	rest := fs.Args()
	if len(rest) != 1 {
		fmt.Fprintln(os.Stderr, "usage: filesort [--dry-run] <rootDir>")
		return 2
	}
	root := rest[0]

	plan, err := sorter.BuildPlan(root, dryRun)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if dryRun {
		// Print a small summary; helpful for future assertions and user feedback.
		fmt.Fprintf(os.Stdout, "dry-run: %d moves planned\n", len(plan.Moves))
		for src, dst := range plan.Moves {
			fmt.Fprintf(os.Stdout, "%s -> %s\n", src, dst)
		}
		return 0
	}

	if err := sorter.Apply(plan); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (int, error) { return len(p), nil }
