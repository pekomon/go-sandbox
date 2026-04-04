# BrickBreaker

BrickBreaker is an Arkanoid-style paddle game scaffold built with Ebiten. This issue intentionally lands the module structure, CI wiring, and deterministic gameplay tests before the real collision loop is implemented.

The current binary opens a placeholder window so the module builds cleanly while the gameplay rules are still being driven by red tests.

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
cd brickbreaker
make deps
make build
```

If you prefer raw Go commands:

```bash
cd brickbreaker
go build -o bin/brickbreaker ./cmd/brickbreaker
```

## Usage

Run the current scaffold:

```bash
cd brickbreaker
./bin/brickbreaker
```

The placeholder window documents the intended controls:
- Left / Right or A / D to move the paddle
- Space to launch the ball
- R to restart after a win or loss
- Esc to quit

## Testing & Coverage

This module is intentionally tests-first at the moment. The gameplay tests describe the expected behavior and currently fail until the follow-up implementation issue lands.

```bash
cd brickbreaker
make test
make cover
```

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Clean exit from the window. |
| `1` | Initialization or runtime failure reported via stderr. |

## Development notes

- Rendering uses [Ebiten v2](https://ebiten.org/) to match the existing desktop games in this repository.
- Core simulation is expected to live in `internal/game`.
- Follow the shared workflow in [../agents.md](../agents.md): tests first, then the feature PR that makes them pass.
