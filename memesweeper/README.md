# MemeSweeper

MemeSweeper is a Minesweeper-inspired puzzler that will use Ebiten to render meme-themed tiles, timers, and flag counters. This directory currently contains the project scaffold so that follow-up issues can add failing tests, the game engine, and the UI loop.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd memesweeper
make deps   # optional; runs go mod tidy
make build  # produces ./bin/memesweeper
```

If you prefer raw Go commands:

```bash
cd memesweeper
go build -o bin/memesweeper ./cmd/memesweeper
```

---

## Usage

The CLI currently exposes a development stub:

```bash
cd memesweeper
./bin/memesweeper --version   # prints 0.1.0-dev
./bin/memesweeper             # placeholder mode (gameplay comes later)
```

Running without flags prints a friendly "not implemented" notice to stderr and returns exit code 1. Issue #63 and beyond will introduce real flags for board sizing, difficulty, and seeded layouts.

### Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Successfully handled the requested action (e.g., `--version`). |
| `1` | Game loop / feature unavailable (current stub behavior). |
| `2` | Invalid CLI usage (flag parsing error). |

### Environment variables

None yet. Future issues will introduce persistence paths and assets directories; document them here once defined.

---

## Testing & coverage

```bash
cd memesweeper
make test   # go test ./...
make cover  # go test ./... -coverprofile=cover.out -covermode=atomic
```

The CI workflow mirrors these targets and uploads `cover.out` artifacts for pull requests touching this module.

---

## Development notes

- Modules stick to the Go standard library by default.
- Gameplay logic will live in `internal/` packages (board generation, adjacency counts, win/loss detection).
- `cmd/memesweeper` owns CLI wiring and the Ebiten event loop once implemented.
- Follow the shared conventions in [../agents.md](../agents.md) for branching, issue flow (tests-first), and PR hygiene.
