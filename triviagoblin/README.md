# TriviaGoblin

TriviaGoblin is a CLI quiz game built around short trivia rounds and deterministic shuffles. It loads questions from a JSON file, optionally filters by category, shuffles deterministically with a seed, and prints a short run summary at the end.

A ready-to-run sample question set is included at `triviagoblin/questions.sample.json`.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd triviagoblin
make deps   # optional; runs go mod tidy
make build  # produces ./bin/triviagoblin
```

If you prefer raw Go commands:

```bash
cd triviagoblin
go build -o bin/triviagoblin ./cmd/triviagoblin
```

---

## Usage

Show the command usage:

```bash
cd triviagoblin
./bin/triviagoblin help
```

Run a quiz using all questions from a JSON file:

```bash
cd triviagoblin
./bin/triviagoblin run --file ./questions.sample.json
```

Run only two questions with a deterministic shuffle:

```bash
cd triviagoblin
./bin/triviagoblin run --file ./questions.sample.json --count 2 --seed 42
```

Filter to a single category:

```bash
cd triviagoblin
./bin/triviagoblin run --file ./questions.sample.json --category geo --seed 7
```

You can also run the CLI without building:

```bash
cd triviagoblin
go run ./cmd/triviagoblin run --file ./questions.sample.json --count 2 --seed 42
```

### Question file format

TriviaGoblin expects a JSON array of objects with:

```json
[
  {
    "prompt": "Capital of France?",
    "answer": "Paris",
    "category": "geo"
  },
  {
    "prompt": "2+2",
    "answer": "4"
  }
]
```

Rules:
- `prompt` is required
- `answer` is required
- `category` is optional
- unknown JSON fields are rejected
- empty files or filters that leave zero questions cause a runtime error
- `questions.sample.json` is a valid example file you can use immediately

---

## Flags

The `run` subcommand supports:

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `--file` | *(required)* | Path to a JSON question file. |
| `--count` | `0` | Number of questions to ask. `0` means ask all loaded questions. |
| `--seed` | `0` | Shuffle seed. Using the same seed gives the same question order. |
| `--category` | *(empty)* | Optional category filter, matched case-insensitively. |

Usage summary:

```text
triviagoblin run --file <path> [--count N] [--seed N] [--category name]
```

---

## Environment variables

No environment variables are defined yet.

---

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Success, including `help` output and completed quiz runs. |
| `1` | Runtime failure, such as file I/O errors, invalid JSON input, or no questions matching the requested category. |
| `2` | Invalid CLI usage, such as unknown commands, invalid flags, missing `--file`, or `--count < 0`. |

Normal output will be printed to stdout; all error messages go to stderr.

---

## Testing & coverage

```bash
cd triviagoblin
make test   # go test ./...
make cover  # go test ./... -coverprofile=cover.out && go tool cover -func cover.out
```

The GitHub Actions workflow mirrors these targets and uploads the `cover.out` artifact for pull requests touching this module.

---

## Development notes

- Keep dependencies stdlib-only unless justified in PR notes.
- The CLI lives in `cmd/triviagoblin` and implementation packages under `internal/`.
- Follow the repository conventions in [../agents.md](../agents.md) for branching strategy, PR templates, and release cadence.
