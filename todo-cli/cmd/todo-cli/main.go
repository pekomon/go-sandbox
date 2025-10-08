package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pekomon/go-sandbox/todo-cli/internal/storage"
	"github.com/pekomon/go-sandbox/todo-cli/internal/tasks"
)

var ErrUnknownCommand = errors.New("unknown command")

type Command struct {
	Name string
	Args []string
}

// ParseCommand inspects the provided CLI arguments (including program name) and
// returns the parsed command. It keeps support for the legacy command parsing
// API used in older tests.
func ParseCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return Command{}, ErrUnknownCommand
	}
	cmd := Command{Name: args[1]}
	rest := args[2:]
	var storagePath string
	filtered := make([]string, 0, len(rest))
	for i := 0; i < len(rest); i++ {
		s := rest[i]
		if s == "--storage" {
			if i+1 >= len(rest) {
				return Command{}, ErrUnknownCommand
			}
			storagePath = rest[i+1]
			i++
			continue
		}
		filtered = append(filtered, s)
	}
	switch cmd.Name {
	case "add":
		if len(filtered) == 0 {
			return Command{}, ErrUnknownCommand
		}
		cmd.Args = []string{strings.Join(filtered, " ")}
	case "list":
		if storagePath != "" {
			cmd.Args = append(cmd.Args, storagePath)
		}
		cmd.Args = append(cmd.Args, filtered...)
	case "done", "rm":
		if len(filtered) != 1 {
			return Command{}, ErrUnknownCommand
		}
		cmd.Args = filtered
		if storagePath != "" {
			cmd.Args = append(cmd.Args, storagePath)
		}
	case "clear":
		if len(filtered) != 0 {
			return Command{}, ErrUnknownCommand
		}
		if storagePath != "" {
			cmd.Args = []string{storagePath}
		}
	default:
		return Command{}, ErrUnknownCommand
	}
	return cmd, nil
}

// Run is separated for testability. Returns exit code.
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: todo-cli <add|list|done|rm|clear> [args]")
		return 2
	}

	// Resolve storage path
	jsonPath, err := storage.DefaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving storage path:", err)
		return 1
	}

	// Commands
	switch args[0] {
	case "add":
		text := strings.TrimSpace(strings.Join(args[1:], " "))
		if text == "" {
			fmt.Fprintln(os.Stderr, "add requires task text")
			return 2
		}
		lock, lerr := storage.AcquireLock(jsonPath)
		if lerr != nil {
			fmt.Fprintln(os.Stderr, lerr)
			return 1
		}
		defer lock.Release()

		list, err := storage.LoadTasks(jsonPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		list = tasks.Add(list, text)
		if err := storage.SaveTasks(jsonPath, list); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Fprintf(os.Stdout, "added #%d\n", list[len(list)-1].ID)
		return 0

	case "list":
		fs := flag.NewFlagSet("list", flag.ContinueOnError)
		reverse := fs.Bool("reverse", false, "reverse order (oldest-first)")
		// prevent flag package from writing to stderr on parse error
		fs.SetOutput(new(nopWriter))
		if err := fs.Parse(args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, "invalid flags")
			return 2
		}
		list, err := storage.LoadTasks(jsonPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if *reverse {
			list = tasks.SortOldestFirst(list)
		} else {
			list = tasks.SortNewestFirst(list)
		}
		for _, t := range list {
			state := " "
			if t.Done {
				state = "x"
			}
			fmt.Fprintf(os.Stdout, "[%s] #%d %s\n", state, t.ID, t.Text)
		}
		return 0

	case "done":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "done requires an ID")
			return 2
		}
		id, convErr := strconv.Atoi(args[1])
		if convErr != nil {
			fmt.Fprintln(os.Stderr, "invalid ID")
			return 2
		}
		lock, lerr := storage.AcquireLock(jsonPath)
		if lerr != nil {
			fmt.Fprintln(os.Stderr, lerr)
			return 1
		}
		defer lock.Release()

		list, err := storage.LoadTasks(jsonPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		list, err = tasks.MarkDone(list, id)
		if err != nil {
			if errors.Is(err, tasks.ErrNotFound) {
				fmt.Fprintln(os.Stderr, "no such task")
				return 2
			}
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := storage.SaveTasks(jsonPath, list); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Fprintf(os.Stdout, "done #%d\n", id)
		return 0

	case "rm":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "rm requires an ID")
			return 2
		}
		id, convErr := strconv.Atoi(args[1])
		if convErr != nil {
			fmt.Fprintln(os.Stderr, "invalid ID")
			return 2
		}
		lock, lerr := storage.AcquireLock(jsonPath)
		if lerr != nil {
			fmt.Fprintln(os.Stderr, lerr)
			return 1
		}
		defer lock.Release()

		list, err := storage.LoadTasks(jsonPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		list, err = tasks.Remove(list, id)
		if err != nil {
			if errors.Is(err, tasks.ErrNotFound) {
				fmt.Fprintln(os.Stderr, "no such task")
				return 2
			}
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := storage.SaveTasks(jsonPath, list); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Fprintf(os.Stdout, "removed #%d\n", id)
		return 0

	case "clear":
		lock, lerr := storage.AcquireLock(jsonPath)
		if lerr != nil {
			fmt.Fprintln(os.Stderr, lerr)
			return 1
		}
		defer lock.Release()

		// Clear regardless of existing content
		if err := storage.SaveTasks(jsonPath, tasks.Clear(nil)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Fprintln(os.Stdout, "cleared")
		return 0

	default:
		fmt.Fprintln(os.Stderr, "unknown command")
		return 2
	}
}

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (int, error) { return len(p), nil }

func main() {
	os.Exit(Run(os.Args[1:]))
}
