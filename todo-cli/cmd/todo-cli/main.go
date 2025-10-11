package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
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
		if os.Getenv("TODO_CLI_MENU") == "1" {
			return runMenu()
		}
		fmt.Fprintln(os.Stderr, "usage: todo-cli <add|list|done|rm|clear> [args]")
		return 2
	}

	if args[0] == "menu" {
		return runMenu()
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

func runMenu() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stdout, "")
		fmt.Fprintln(os.Stdout, "TODO CLI MENU")
		fmt.Fprintln(os.Stdout, "----------------")
		fmt.Fprintln(os.Stdout, "1) Add task")
		fmt.Fprintln(os.Stdout, "2) List tasks")
		fmt.Fprintln(os.Stdout, "3) Mark done")
		fmt.Fprintln(os.Stdout, "4) Remove task")
		fmt.Fprintln(os.Stdout, "5) Clear tasks")
		fmt.Fprintln(os.Stdout, "0) Exit")
		fmt.Fprint(os.Stdout, "Select an option: ")

		choice, err := readTrimmedLine(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stdout)
				return 0
			}
			fmt.Fprintln(os.Stderr, "input error:", err)
			return 1
		}

		switch choice {
		case "1":
			text, err := promptLine(reader, "Enter task description: ")
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Fprintln(os.Stdout)
					return 0
				}
				fmt.Fprintln(os.Stderr, "input error:", err)
				return 1
			}
			if text == "" {
				fmt.Fprintln(os.Stdout, "no text entered")
				continue
			}
			if exit := Run([]string{"add", text}); exit == 1 {
				return 1
			}
		case "2":
			if exit := menuList(); exit != -1 {
				return exit
			}
		case "3":
			id, ok, exitCode := promptID(reader, "mark done")
			if exitCode >= 0 {
				return exitCode
			}
			if !ok {
				continue
			}
			if exit := Run([]string{"done", id}); exit == 1 {
				return 1
			}
		case "4":
			id, ok, exitCode := promptID(reader, "remove")
			if exitCode >= 0 {
				return exitCode
			}
			if !ok {
				continue
			}
			if exit := Run([]string{"rm", id}); exit == 1 {
				return 1
			}
		case "5":
			if exit := Run([]string{"clear"}); exit == 1 {
				return 1
			}
		case "0":
			fmt.Fprintln(os.Stdout, "Goodbye!")
			return 0
		default:
			if strings.EqualFold(choice, "exit") {
				fmt.Fprintln(os.Stdout, "Goodbye!")
				return 0
			}
			fmt.Fprintln(os.Stderr, "invalid selection")
		}
	}
}

func promptLine(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)
	return readTrimmedLine(reader)
}

func promptID(reader *bufio.Reader, action string) (string, bool, int) {
	fmt.Fprintf(os.Stdout, "Enter task ID to %s: ", action)
	line, err := readTrimmedLine(reader)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stdout)
			return "", false, 0
		}
		fmt.Fprintln(os.Stderr, "input error:", err)
		return "", false, 1
	}
	if line == "" {
		fmt.Fprintln(os.Stdout, "no ID entered")
		return "", false, -1
	}
	if _, convErr := strconv.Atoi(line); convErr != nil {
		fmt.Fprintln(os.Stdout, "invalid ID")
		return "", false, -1
	}
	return line, true, -1
}

func menuList() int {
	jsonPath, err := storage.DefaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving storage path:", err)
		return 1
	}
	list, err := storage.LoadTasks(jsonPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if len(list) == 0 {
		fmt.Fprintln(os.Stdout, "No tasks found.")
		return -1
	}
	list = tasks.SortNewestFirst(list)
	for _, t := range list {
		state := " "
		if t.Done {
			state = "x"
		}
		fmt.Fprintf(os.Stdout, "[%s] #%d %s\n", state, t.ID, t.Text)
	}
	return -1
}

func readTrimmedLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			if len(line) == 0 {
				return "", io.EOF
			}
			return strings.TrimSpace(line), nil
		}
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func main() {
	os.Exit(Run(os.Args[1:]))
}
