package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pekomon/go-sandbox/dungeondice/internal/dungeondice"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stderr)
		return 2
	}

	switch args[0] {
	case "run":
		return runRun(args[1:])
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printUsage(os.Stderr)
		return 2
	}
}

func runRun(args []string) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(new(nopWriter))
	var (
		className string
		seed      int64
		rooms     int
	)
	fs.StringVar(&className, "class", "", "player class")
	fs.IntVar(&rooms, "rooms", 3, "number of rooms in the run")
	fs.Int64Var(&seed, "seed", 0, "deterministic seed (0 = random)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "invalid flags")
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected arguments")
		return 2
	}
	if strings.TrimSpace(className) == "" {
		fmt.Fprintln(os.Stderr, "missing --class")
		return 2
	}
	if rooms <= 0 {
		fmt.Fprintln(os.Stderr, "rooms must be positive")
		return 2
	}

	summary, err := dungeondice.SimulateRun(dungeondice.RunConfig{
		Class: className,
		Seed:  seed,
		Rooms: rooms,
	})
	if err != nil {
		if errors.Is(err, dungeondice.ErrUnknownClass) {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintf(os.Stderr, "valid classes: %s\n", strings.Join(dungeondice.ClassNames(), ", "))
			return 2
		}
		if errors.Is(err, dungeondice.ErrInvalidRooms) {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	printSummary(os.Stdout, summary)
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: dungeondice run --class <name> [--rooms N] [--seed N]")
}

func printSummary(w io.Writer, summary dungeondice.RunSummary) {
	fmt.Fprintln(w, "Run summary")
	fmt.Fprintf(w, "Class: %s\n", summary.Class)
	fmt.Fprintf(w, "Seed: %d\n", summary.Seed)
	fmt.Fprintf(w, "Rooms: %d\n", summary.Rooms)
	fmt.Fprintf(w, "Cleared: %d\n", summary.Cleared)
	fmt.Fprintf(w, "State: %s\n", summary.State)
	fmt.Fprintf(w, "Final HP: %d/%d\n", summary.FinalHP, summary.FinalMaxHP)
	fmt.Fprintf(w, "Rounds: %d\n", summary.Rounds)
}

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (int, error) { return len(p), nil }
