# Speedpads

Speedpads is a Speden Spelit inspired memory-speed game built with Ebiten. The current implementation includes the core show/input loop, deterministic sequence generation, round-to-round speed ramping, and watchdog timeouts for stalled input.

## Installation

Requirements: Go 1.25+, `make`, and a desktop graphics environment for Ebiten.

```bash
cd speedpads
make build
```

Or build directly:

```bash
cd speedpads
go build -o bin/speedpads ./cmd/speedpads
```

## Usage

Run the game:

```bash
cd speedpads
./bin/speedpads
```

Or run without building:

```bash
cd speedpads
go run ./cmd/speedpads
```

Controls:
- `Space` to start
- `D`, `F`, `J`, `K` for the four pads
- `R` to restart after game over
- `Esc` to quit

## Testing & Coverage

The gameplay tests cover the deterministic sequence, show phase, input phase, watchdog timeout, and round progression behavior.

```bash
cd speedpads
make test
make cover
```

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Clean exit from the window. |
| `1` | Initialization or runtime failure reported via stderr. |

## Development notes

- Rendering uses [Ebiten v2](https://ebiten.org/) to match the existing GUI games in this repository.
- Core game rules live in `internal/game`.
- Follow the shared repository conventions in [../agents.md](../agents.md).
