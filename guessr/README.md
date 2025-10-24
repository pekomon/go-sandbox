# guessr

Number guessing CLI (stdlib-only). The program picks a secret number and you guess with hints **higher/lower** until you get it right or run out of attempts.

## Features
- Flags: `--max`, `--attempts`, `--seed` (deterministic RNG for reproducible runs)
- Hints: **higher** / **lower**, success message **Correct!**
- JSON-backed stats (games played, wins, average guesses)
- Pure Go standard library; offline by default

## Installation

Build from source (this submodule lives in `guessr/`):

```bash
# From repo root
cd guessr
make build
# Binary at: ./bin/guessr
```

Or with plain Go:

```bash
cd guessr
go build -o bin/guessr ./cmd/guessr
```

## Usage

Basic run (defaults: `--max=100`, `--attempts=7`, random seed):

```bash
./bin/guessr
# Type one integer per line and press Enter each time, e.g.:
# 50
# 75
# ...
```

Deterministic run (useful for demos/tests):

```bash
./bin/guessr --max=10 --attempts=5 --seed=42
```

You can script a quick interaction:

```bash
# Example: feed a few guesses from stdin
printf "5\n7\n3\n" | ./bin/guessr --max=10 --attempts=5 --seed=42
```

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--max` | `100` | Upper bound of the secret number (range is 1..max). |
| `--attempts` | `7` | Maximum number of guesses. |
| `--seed` | `0` | RNG seed; `0` = non-deterministic (time-based), non-zero = fixed. |

## Output & Behavior

On each guess: prints **higher** or **lower**.

On success: prints **Correct!** and exits 0.

On attempts exhausted: prints **Game over! The number was N.** and exits 0.

Non-integer input: prints **Invalid input** and does not consume an attempt.

## Stats (JSON persistence)

By default, stats are stored at:

`~/.guessr/stats.json`

Override the path (useful for local experiments or tests):

```bash
export GUESSR_STATS_PATH="$(mktemp -d)/stats.json"
./bin/guessr --max=10 --attempts=3 --seed=1
```

The JSON schema is minimal:

```json
{
  "games": 3,
  "wins": 2,
  "total_guesses": 9,
  "average_guesses": 3.0
}
```

`average_guesses = total_guesses / games`.

## Testing & Coverage

Run tests for this subproject:

```bash
cd guessr
make test
```

Generate a local coverage report:

```bash
cd guessr
make cover
# Produces cover.out and a summary
```

GitHub Actions runs tests on PRs affecting `guessr/` and uploads coverage artifacts.

## Exit codes

- `0` success (win or clean game over)
- `1` runtime/storage error (e.g., I/O failure)
- `2` usage error (flag parse error)

Errors are printed to stderr; normal game output goes to stdout.

## Troubleshooting

### “Invalid input” appears
The program expects integers per line. The message does not consume attempts; just enter a valid number.

### Stats file is corrupted
The program returns an error on JSON parsing. Move the file aside and re-run:

```bash
mv ~/.guessr/stats.json ~/.guessr/stats.json.bak.$(date +%s)
```

### Different results between runs
Set a non-zero `--seed` to make runs deterministic.
