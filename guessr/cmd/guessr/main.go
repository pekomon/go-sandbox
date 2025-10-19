package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pekomon/go-sandbox/guessr/internal/game"
	"github.com/pekomon/go-sandbox/guessr/internal/stats"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	var (
		max      = flag.NewFlagSet("guessr", flag.ContinueOnError)
		maxVal   int
		attempts int
		seed     int64
	)
	max.IntVar(&maxVal, "max", 100, "upper bound for the secret number (inclusive range 1..max)")
	max.IntVar(&attempts, "attempts", 7, "max number of guesses")
	max.Int64Var(&seed, "seed", 0, "deterministic seed (0 = random)")

	// silence default usage on parse error; print minimal hint instead
	max.SetOutput(new(nopWriter))
	if err := max.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "invalid flags")
		return 2
	}

	store, err := stats.NewStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	opts := game.Options{Max: maxVal, Attempts: attempts, Seed: seed}
	if err := game.Run(os.Stdin, os.Stdout, opts, store); err != nil {
		// Run should not normally return an error (it handles input itself),
		// but if it does, treat as runtime error.
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (int, error) { return len(p), nil }
