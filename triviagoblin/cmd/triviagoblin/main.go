package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pekomon/go-sandbox/triviagoblin/internal/quiz"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) == 0 {
		printUsage(errOut)
		return 2
	}

	switch args[0] {
	case "run":
		return runQuiz(args[1:], in, out, errOut)
	case "help", "-h", "--help":
		printUsage(out)
		return 0
	default:
		fmt.Fprintf(errOut, "unknown command: %s\n", args[0])
		printUsage(errOut)
		return 2
	}
}

func runQuiz(args []string, in io.Reader, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		filePath string
		category string
		count    int
		seed     int64
	)
	fs.StringVar(&filePath, "file", "", "path to a JSON question file")
	fs.IntVar(&count, "count", 0, "number of questions to ask (0 = all)")
	fs.Int64Var(&seed, "seed", 0, "shuffle seed")
	fs.StringVar(&category, "category", "", "optional category filter")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(errOut, "invalid flags")
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(errOut, "unexpected arguments")
		return 2
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		fmt.Fprintln(errOut, "missing --file")
		return 2
	}
	if count < 0 {
		fmt.Fprintln(errOut, "count must be >= 0")
		return 2
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	defer file.Close()

	questions, err := quiz.LoadQuestions(file)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}

	questions = quiz.FilterByCategory(questions, category)
	if len(questions) == 0 {
		if strings.TrimSpace(category) != "" {
			fmt.Fprintf(errOut, "no questions matched category %q\n", category)
		} else {
			fmt.Fprintln(errOut, "no questions available")
		}
		return 1
	}

	summary, err := quiz.Run(in, out, questions, quiz.Config{
		Seed:  seed,
		Count: count,
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}

	printSummary(out, summary, seed, category)
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: triviagoblin run --file <path> [--count N] [--seed N] [--category name]")
}

func printSummary(w io.Writer, summary quiz.Summary, seed int64, category string) {
	category = strings.TrimSpace(category)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run summary")
	fmt.Fprintf(w, "Questions available: %d\n", summary.Total)
	fmt.Fprintf(w, "Questions asked: %d\n", summary.Asked)
	fmt.Fprintf(w, "Correct answers: %d\n", summary.Correct)
	fmt.Fprintf(w, "Seed: %d\n", seed)
	if category != "" {
		fmt.Fprintf(w, "Category: %s\n", category)
	}
}
