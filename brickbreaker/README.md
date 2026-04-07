# BrickBreaker

BrickBreaker is an Arkanoid-style paddle game built with Ebiten. It runs fully offline, opens a desktop window, and keeps the core simulation deterministic under `internal/game` so the gameplay rules stay testable.

The current version includes three built-in levels, a level-by-level speed ramp, score/lives HUD, and quick restarts after a win or loss.

## Requirements

- Go 1.25 or newer
- macOS, Windows, or Linux with a working desktop graphics environment

Headless environments are not supported for gameplay, though the unit tests still run offline.

## Installation

Build with the module Makefile:

```bash
cd brickbreaker
make build
```

The binary will be written to `./bin/brickbreaker`.

Or build directly with Go:

```bash
cd brickbreaker
go build -o bin/brickbreaker ./cmd/brickbreaker
```

Ebiten is fetched automatically by Go modules; no extra runtime setup is required beyond normal desktop graphics support.

## Usage

Run the built game:

```bash
cd brickbreaker
./bin/brickbreaker
```

Or run it without building first:

```bash
cd brickbreaker
go run ./cmd/brickbreaker
```

The game starts in serve state with the ball attached to the paddle. Launch with `Space`, clear each level’s bricks, and survive long enough to reach the final win screen.

## Controls

| Key(s) | Action |
| ---- | ------- |
| Left / Right | Move the paddle |
| A / D | Move the paddle |
| Space | Launch the ball |
| R | Restart after win/loss |
| Space / Enter | Restart after win/loss |
| Esc | Quit |

## Gameplay

- The game ships with three deterministic built-in levels.
- Clearing every brick on a level loads the next layout automatically.
- Ball speed increases on later levels.
- Dropping the ball costs one life and returns to serve state.
- Losing all lives ends the run.
- Clearing the final level shows the win state without closing the window.

The HUD shows:
- current level
- remaining lives
- score
- ball speed while the ball is in play

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Clean exit from the window. |
| `1` | Initialization or runtime failure reported via stderr. |

## Testing & Coverage

Run the unit tests:

```bash
cd brickbreaker
make test
```

Generate local coverage:

```bash
cd brickbreaker
make cover
```

The test suite covers the deterministic simulation in `internal/game`, including:
- serve-state movement
- launch behavior
- wall, paddle, and brick collisions
- life loss and game-over handling
- level progression and speed ramp behavior

## Development notes

- The desktop UI lives in `cmd/brickbreaker`.
- Core gameplay rules live in `internal/game`.
- Rendering uses [Ebiten v2](https://ebiten.org/) to match the other GUI games in this repository.
- Follow the shared repository conventions in [../agents.md](../agents.md).

## Troubleshooting

- If the game window does not open, check that your desktop graphics environment is working normally.
- macOS may show Ebiten-related `CVDisplayLink` deprecation warnings during local builds; these warnings are harmless for current development.
