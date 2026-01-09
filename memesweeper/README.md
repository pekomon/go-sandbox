# MemeSweeper

MemeSweeper is a Minesweeper-inspired puzzler rendered with Ebiten, featuring meme-themed tiles, flags, and a basic win/loss loop.

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

The CLI launches the desktop UI:

```bash
cd memesweeper
./bin/memesweeper --version   # prints 0.1.0-dev
./bin/memesweeper             # starts the Ebiten window
./bin/memesweeper --difficulty hard
```

You can also run without building:

```bash
cd memesweeper
go run ./cmd/memesweeper
```

Controls:

- Left click: reveal tile
- Right click: toggle flag
- 1 / 2 / 3: easy / medium / hard
- R / Space / Enter: restart (keeps current difficulty)
- Esc: quit

Issues #67-68 will introduce difficulty presets and CLI/UI selection.

### Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Successfully handled the requested action (e.g., `--version`). |
| `1` | Game loop failure or asset loading error. |
| `2` | Invalid CLI usage (flag parsing error). |

### Environment variables

None yet. Assets are loaded from `memesweeper/assets/`.

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
- Gameplay logic lives in `internal/board`, UI logic in `internal/ui`.
- `cmd/memesweeper` owns CLI wiring and the Ebiten event loop.
- Follow the shared conventions in [../agents.md](../agents.md) for branching, issue flow (tests-first), and PR hygiene.
