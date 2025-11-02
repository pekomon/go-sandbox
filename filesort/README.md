# filesort

CLI tool to classify files in a directory into subfolders by type (based on extension).  
Implemented in pure Go (stdlib only).

## Features

- Classifies files into:
  - `images/` — `.jpg`, `.jpeg`, `.png`, `.gif`
  - `docs/` — `.pdf`, `.doc`, `.docx`, `.txt`, `.md`
  - `videos/` — `.mp4`, `.mov`, `.avi`
  - `other/` — everything else
- Supports **dry-run mode** (`--dry-run`) to preview planned moves without modifying files.
- Non-recursive for simplicity; acts only on the top-level of the given directory.

## Installation

```bash
# From repo root
cd filesort
make build
# Binary available at: ./bin/filesort
```

Or directly with Go:

```bash
cd filesort
go build -o bin/filesort ./cmd/filesort
```

## Usage

Sort files inside a directory:

```bash
./bin/filesort ~/Downloads
```

Show planned actions without moving files:

```bash
./bin/filesort --dry-run ~/Downloads
```

Example output:

```text
dry-run: 5 moves planned
/home/user/Downloads/photo.JPG -> /home/user/Downloads/images/photo.JPG
/home/user/Downloads/notes.md -> /home/user/Downloads/docs/notes.md
...
```

After a non–dry-run execution, you’ll see:

```text
$ tree ~/Downloads
Downloads/
├── docs/
│   └── notes.md
├── images/
│   └── photo.JPG
├── videos/
│   └── clip.MP4
└── other/
    └── archive.tar.gz
```

## Flags

| Flag | Description |
| ---- | ----------- |
| `--dry-run` | Compute and display the plan without moving files. |

## Exit codes

- `0` — Success (plan printed in dry-run mode or moves applied without errors)
- `1` — Runtime failure (I/O issues, invalid destination plan, move failure)
- `2` — Usage error (flag parse failure or missing/extra arguments)

Errors are printed to stderr; dry-run and progress messages go to stdout.

## Testing & Coverage

```bash
cd filesort
make test
make cover
```

These run unit tests against the core planner and filesystem operations in a temporary directory.  
The tests are hermetic and do not modify your actual filesystem.

GitHub Actions automatically executes these on PRs touching `filesort/` and uploads coverage artifacts.

## Implementation notes

- Uses only the Go standard library.
- Non-recursive: only top-level files are processed.
- All moves use atomic `os.Rename`.
- Destination directories are created on demand with `os.MkdirAll`.

## Development conventions

Follow the shared guidelines in [../agents.md](../agents.md). That document defines:

- branching & commit style (tests-first, one feature per PR),
- Go toolchain and CI setup,
- rules for adding dependencies and writing offline tests.

## Troubleshooting

- **“Permission denied”** — ensure you have write permission to the target directory.
- **“not a directory”** — pass a valid directory path (not a file).
- **CI failure after merge** — run `make ensure-tidy` to verify your `go.mod` and `go.sum` are consistent.
