package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pekomon/go-sandbox/memesweeper/internal/ui"
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

	cfg := ui.Config{
		Rows:      10,
		Cols:      10,
		MemeCount: 15,
		TileSize:  32,
	}
	if err := ui.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
