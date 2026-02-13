# DungeonDice

DungeonDice is a CLI roguelike simulation built around dice-driven combat rounds and class-based runs. The `run` command simulates a full run and prints a summary.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd dungeondice
make deps   # optional; runs go mod tidy
make build  # produces ./bin/dungeondice
```

If you prefer raw Go commands:

```bash
cd dungeondice
go build -o bin/dungeondice ./cmd/dungeondice
```

---

## Usage

Build the binary first, then run a simulated adventure:

```bash
cd dungeondice
./bin/dungeondice run --class warrior --rooms 3 --seed 42
```

Try a different class or let the seed be random:

```bash
cd dungeondice
./bin/dungeondice run --class rogue --rooms 4
./bin/dungeondice run --class mystic --rooms 2 --seed 7
```

Available classes: `warrior`, `rogue`, `mystic`.

Output looks like:

```text
Run summary
Class: Warrior
Seed: 42
Rooms: 3
Cleared: <n>
State: <state>
Final HP: <n>/<n>
Rounds: <n>
```

---

## Flags

The `run` subcommand accepts:

- `--class` (required): class name (`warrior`, `rogue`, `mystic`)
- `--rooms` (optional): number of rooms in the run (default `3`, must be > 0)
- `--seed` (optional): deterministic seed (`0` = random)

---

## Environment variables

No environment variables are defined yet.

---

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Success. |
| `1` | Runtime failure. |
| `2` | Invalid CLI usage. |

Normal output will be printed to stdout; all error messages go to stderr.

---

## Testing & coverage

```bash
cd dungeondice
make test   # go test ./...
make cover  # go test ./... -coverprofile=cover.out && go tool cover -func cover.out
```

The GitHub Actions workflow mirrors these targets and uploads the `cover.out` artifact for pull requests touching this module.

---

## Development notes

- Keep dependencies stdlib-only unless justified in PR notes.
- The CLI lives in `cmd/dungeondice` and implementation packages under `internal/`.
- Follow the repository conventions in [../agents.md](../agents.md) for branching strategy, PR templates, and release cadence.
