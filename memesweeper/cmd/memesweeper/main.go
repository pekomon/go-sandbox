package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pekomon/go-sandbox/memesweeper/internal/board"
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
	var difficulty string
	fs.BoolVar(&showVersion, "version", false, "print the MemeSweeper version")
	fs.StringVar(&difficulty, "difficulty", "medium", "difficulty preset: easy, medium, hard")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if showVersion {
		fmt.Fprintln(os.Stdout, version)
		return 0
	}

	preset := board.Preset(difficulty)
	if _, err := board.PresetConfig(preset, 0); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	cfg := ui.Config{
		Preset:   preset,
		TileSize: 32,
	}
	if err := ui.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
