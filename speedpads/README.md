# Speedpads

Speedpads is a Speden Spelit inspired memory-speed game built with Ebiten. The current version includes a playable four-pad loop, deterministic sequence generation, playback audio cues, a game-over tone, round-to-round speed ramping, and watchdog timeouts for stalled input.

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
- `Space` or `Enter` to start
- `D`, `F`, `J`, `K` for the four pads
- `R` to restart after game over
- `Esc` to quit

Gameplay notes:
- The game replays the full sequence each round and then waits for your input.
- If you already know the sequence, you can start pressing during playback. The first press cuts the remaining playback and switches immediately to input mode.
- The game plays distinct tones when it flashes sequence pads and a separate descending tone on game over.
- Your own button presses stay silent so only the playback and failure cues make sound.

## Testing & Coverage

The gameplay tests cover deterministic sequence generation, show-to-input transitions, early-input handling, watchdog timeout, and round progression behavior.

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
- Audio cues are synthesized in code with Ebiten's audio package; there are no external sound assets.
- Core game rules live in `internal/game`.
- Follow the shared repository conventions in [../agents.md](../agents.md).
