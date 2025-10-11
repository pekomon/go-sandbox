# todo-cli

Local TODO manager (CLI). Stores tasks in a JSON file on your machine.

## Features
- Commands: `add`, `list`, `done <id>`, `rm <id>`, `clear`
- Interactive menu: `todo-cli menu` or `TODO_CLI_MENU=1 todo-cli`
- Default sort: newest first (by ID, descending); `list --reverse` flips order
- JSON persistence at `~/.todo-cli/tasks.json` (overridable)
- Standard library only; offline by default

## Installation

Build from source (module lives in this subfolder):

```bash
# From repo root
cd todo-cli
make build
# Binary at: ./bin/todo-cli
Or with plain go:

bash
Copy code
cd todo-cli
go build -o bin/todo-cli ./cmd/todo-cli
Quick Start
bash
Copy code
# Add a couple of tasks
./bin/todo-cli add "Read Go docs"
./bin/todo-cli add "Write unit tests"

# List tasks (newest first)
./bin/todo-cli list
# Example output:
# [ ] #2 Write unit tests
# [ ] #1 Read Go docs

# Mark a task done
./bin/todo-cli done 1
# stdout: done #1

# List again
./bin/todo-cli list
# [ ] #2 Write unit tests
# [x] #1 Read Go docs

# Remove a task
./bin/todo-cli rm 2
# stdout: removed #2

# Clear all tasks
./bin/todo-cli clear
# stdout: cleared
Command Reference
bash
Copy code
todo-cli add <text...>          # Add a new task with the given text
todo-cli list [--reverse]       # List tasks (newest first; --reverse = oldest first)
todo-cli done <id>              # Mark the task with ID as done
todo-cli rm <id>                # Remove the task with ID
todo-cli clear                  # Remove all tasks
todo-cli menu                   # Launch the interactive text menu
Exit codes
0 success

1 runtime/storage error (e.g., I/O, corrupted JSON, lock acquisition fails)

2 usage error (unknown command, missing/invalid arguments)

Errors are printed to stderr; normal output to stdout.

Persistence
By default, tasks are stored at:

Path: ~/.todo-cli/tasks.json

You can override the location for local runs or tests:

Environment variable: TODO_CLI_PATH=/path/to/tasks.json

Example:

bash
Copy code
export TODO_CLI_PATH="$(mktemp -d)/tasks.json"
./bin/todo-cli add "Temporary task"
./bin/todo-cli list
Locking & concurrency
To reduce the chance of concurrent writes, the CLI uses a best-effort lock file next to tasks.json:

Lock file: tasks.lock

Behavior: creates the file with O_CREATE|O_EXCL; if it exists, the command fails

Limitation: on crash/kill, a stale lock may remain; remove it manually if needed

This is intentionally simple. If stronger guarantees are required later, we can upgrade the locking strategy in a future PR.

Testing
Run tests for this subproject:

bash
Copy code
cd todo-cli
make test
Generate coverage locally:

bash
Copy code
cd todo-cli
make cover
# Outputs cover.out and a summary
CI runs go test ./... -coverprofile=cover.out on every PR/push and uploads artifacts.

Development Notes
IDs are sequential starting at 1 (new tasks get max(id)+1)

“Newest first” is implemented via ID ordering to keep the data model small

JSON schema (minimal):

json
Copy code
[
  { "id": 1, "text": "example", "done": false }
]
Troubleshooting
"another process may be running (lock exists)"
A previous run created tasks.lock. If you are sure no other instance runs, remove the lock file:

bash
Copy code
rm -f ~/.todo-cli/tasks.lock
Or point to a different JSON path using TODO_CLI_PATH.

"invalid ID" / "no such task"
Use todo-cli list to inspect current IDs.

Corrupted JSON
The CLI returns an error instead of silently resetting. Fix or move the file, or start with a fresh path:

bash
Copy code
mv ~/.todo-cli/tasks.json ~/.todo-cli/tasks.json.bak.$(date +%s)
