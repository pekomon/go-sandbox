# todo-cli

Local todo manager for the terminal. Tasks live in a JSON file on your machine, and the app works fully offline.

## Features
- Straightforward commands: `add`, `list`, `done <id>`, `rm <id>`, `clear`, `menu`
- Interactive menu with arrow-key navigation (surveys) and a text-mode fallback
- Tasks stored at `~/.todo-cli/tasks.json` (configurable via env var)
- Default ordering shows the newest entries first; `--reverse` lists oldest first
- Standard-library dependencies only

## Installation

Build from source (module is in this directory).

```bash
# From the repo root
cd todo-cli
make build
# Binary is written to ./bin/todo-cli
```

Or build directly with Go:

```bash
cd todo-cli
go build -o bin/todo-cli ./cmd/todo-cli
```

## Usage

### Quick start

```bash
# Add tasks
./bin/todo-cli add "Read Go docs"
./bin/todo-cli add "Write unit tests"

# List tasks (newest first by default)
./bin/todo-cli list

# Mark a task complete
./bin/todo-cli done 1

# Remove a task
./bin/todo-cli rm 2

# Clear all tasks
./bin/todo-cli clear
```

### Command reference

```text
todo-cli add <text...>        Add a new task with the provided text
todo-cli list [--reverse]     List tasks (newest first; --reverse flips to oldest first)
todo-cli done <id>            Mark the task with ID as done
todo-cli rm <id>              Remove the task with ID
todo-cli clear                Remove all tasks
todo-cli menu                 Launch the interactive menu UI
```

### Exit codes

| Code | Meaning                                                                 |
|------|-------------------------------------------------------------------------|
| 0    | Success                                                                 |
| 1    | Runtime/storage error (I/O failure, corrupted JSON, lock acquisition)   |
| 2    | Usage error (unknown command, missing/invalid arguments)                |

Errors are written to stderr; normal output uses stdout.

## Interactive menu

Launch it with `todo-cli menu`, or run `TODO_CLI_MENU=1 todo-cli` to enter the menu automatically when no arguments are provided. When the app detects an interactive terminal it uses arrow keys (via [survey.Select](https://github.com/AlecAivazis/survey)) and `Enter` to pick an action:

```
TODO CLI MENU
> Add task
  List tasks
  Mark done
  Remove task
  Clear tasks
  Exit
```

- Press `Enter` to confirm a highlighted item.
- Choose **Exit** or press `Ctrl+C` to quit; either path returns exit code 0.
- If the process cannot use the survey UI (for example in a non-TTY shell), it falls back to a numbered text menu:
  ```
  TODO CLI MENU
  ----------------
  1) Add task
  2) List tasks
  3) Mark done
  4) Remove task
  5) Clear tasks
  0) Exit
  ```
  Type the number (or `exit`) and press `Enter` to proceed.

IDs requested by menu prompts must be numeric; blank or invalid input keeps the menu open without running a command.

## Configuration

- `TODO_CLI_PATH` overrides the default JSON location.
  ```bash
  export TODO_CLI_PATH="$(mktemp -d)/tasks.json"
  ./bin/todo-cli add "Temporary task"
  ./bin/todo-cli list
  ```
- `TODO_CLI_MENU=1` makes `todo-cli` (with no arguments) launch directly into the menu.

### Persistence & locking

- Default data path: `~/.todo-cli/tasks.json`
- Lock file: `~/.todo-cli/tasks.lock`
- The CLI writes a temporary `.tmp` file and renames it for atomic saves.
- Locking uses a simple `O_CREATE|O_EXCL` file. If another process holds the lock, commands exit with code 1 and report `another process may be running (lock exists)`.
- After an unexpected crash you may need to remove a stale lock manually:
  ```bash
  rm -f ~/.todo-cli/tasks.lock
  ```

## Testing

```bash
cd todo-cli
make test
```

Generate coverage locally:

```bash
cd todo-cli
make cover
# Produces cover.out and a summary
```

CI runs `go test ./... -coverprofile=cover.out` on every PR/push and stores artifacts.

## Development notes

- IDs are sequential starting at 1; new tasks get `max(id) + 1`.
- Newest-first listing is implemented via ID ordering to avoid storing timestamps.
- Minimal JSON structure:
  ```json
  [
    { "id": 1, "text": "example task", "done": false }
  ]
  ```

## Troubleshooting

- **"another process may be running (lock exists)"**  
  Ensure no other `todo-cli` instance is mutating the same file. Remove the lock with `rm -f ~/.todo-cli/tasks.lock`, or point `TODO_CLI_PATH` at a different file.

- **"invalid ID" or "no such task"**  
  Run `todo-cli list` to check the current task IDs, then retry.

- **Corrupted JSON**  
  The CLI refuses to overwrite invalid data. Fix the file manually or back it up and start fresh:
  ```bash
  mv ~/.todo-cli/tasks.json ~/.todo-cli/tasks.json.bak.$(date +%s)
  ```
