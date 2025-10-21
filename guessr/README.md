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


Or with plain Go:

cd guessr
go build -o bin/guessr ./cmd/guessr

Usage

Basic run (defaults: --max=100, --attempts=7, random seed):

./bin/guessr
# Type one integer per line and press Enter each time, e.g.:
# 50
# 75
# ...


Deterministic run (useful for demos/tests):

./bin/guessr --max=10 --attempts=5 --seed=42


You can script a quick interaction:

# Example: feed a few guesses from stdin
printf "5\n7\n3\n" | ./bin/guessr --max=10 --attempts=5 --seed=42

Flags
FlagDefaultDescription
--max100Upper bound of the secret number (range is 1..max).
--attempts7Maximum number of guesses.
--seed0RNG seed; 0 = non-deterministic (time-based). Non-zero = fixed.
Output & Behavior

On each guess: prints higher or lower.

On success: prints Correct! and exits 0.

On attempts exhausted: prints Game over! The number was N. and exits 0.

Non-integer input: prints Invalid input and does not consume an attempt.

Stats (JSON persistence)

By default, stats are stored at:

Path: ~/.guessr/stats.json

Override the path (useful for local experiments or tests):

export GUESSR_STATS_PATH="$(mktemp -d)/stats.json"
./bin/guessr --max=10 --attempts=3 --seed=1


The JSON schema is minimal:

{
  "games": 3,
  "wins": 2,
  "total_guesses": 9,
  "average_guesses": 3.0
}


average_guesses = total_guesses / games.

Testing & Coverage

Run tests for this subproject:

cd guessr
make test


Generate a local coverage report:

cd guessr
make cover
# Produces cover.out and a summary


GitHub Actions runs tests on PRs affecting guessr/ and uploads coverage artifacts.

Exit codes

0 success (win or clean game over)

1 runtime/storage error (e.g., I/O failure)

2 usage error (flag parse error)

Errors are printed to stderr; normal game output goes to stdout.

Troubleshooting

“Invalid input” appears
The program expects integers per line. The message does not consume attempts; just enter a valid number.

Stats file is corrupted
The program returns an error on JSON parsing. Move the file aside and re-run:

mv ~/.guessr/stats.json ~/.guessr/stats.json.bak.$(date +%s)


Different results between runs
Set a non-zero --seed to make runs deterministic.
