# snake

Arcade-style snake clone built with Ebiten. The binary opens a desktop window, runs fully offline, and lets you restart instantly after each run without quitting.

## Features
- Native window rendered via Ebiten; no browser or external runtime required.
- Arrow keys or WASD for steering with buffered turns (no 180° reversals).
- Press `R` to restart immediately; `Space` or `Enter` restarts after game over.
- Overlay shows live score and movement speed; the snake accelerates every five apples.
- Game logic lives in `internal/game` so unit tests stay deterministic.

## Requirements
- Go 1.25 or newer
- macOS, Windows, or Linux with X11/Wayland and functional GPU drivers (headless setups are unsupported)

## Installation

Build with the project Makefile (writes binaries to `./bin/`):

```bash
cd snake
make build
# Binary: ./bin/snake
```

Or build directly with Go:

```bash
cd snake
go build -o bin/snake ./cmd/snake
```

Modules fetch Ebiten automatically; no manual dependency install is needed.

## Usage

Launch the game from the subproject root:

```bash
./bin/snake
```

The window defaults to 24×18 tiles at 24 px each (576×432). Quit with `Esc` at any time.

### Controls

| Key(s)           | Action                                  |
|------------------|-----------------------------------------|
| Arrow keys / WASD| Turn the snake (no instant reversal)    |
| R                | Restart immediately                     |
| Space / Enter    | Restart after game over                 |
| Esc              | Exit the window                         |

### Gameplay flow
- Collect apples to grow the snake and add to your score.
- Colliding with walls or yourself ends the run but keeps the window open for quick restarts.
- Speed ramps up gradually as your score increases, capped at a 2-frame movement delay.

## Exit codes
- `0` — Clean exit (quit or game over screen closed)
- `1` — Unexpected initialization/runtime failure (reported via stderr)

## Testing & Coverage

Run the unit tests:

```bash
cd snake
make test      # go test ./...
```

Generate coverage locally:

```bash
cd snake
make cover     # go test ./... -coverprofile=cover.out
```

## Development

- Rendering + input depends on [Ebiten v2](https://ebiten.org/); all other logic is Go stdlib.
- Follow the shared contributor flow documented in [../agents.md](../agents.md).

## Troubleshooting

- **“command not found: ebiten” during build** — always build with `go build`; modules pull Ebiten automatically.
- **Window fails to open (Linux)** — install required X11/Wayland libraries and GPU drivers. Remote/headless environments are unsupported.
- **macOS CVDisplayLink warnings** — safe to ignore on macOS 15 until upstream Ebiten releases an update.
