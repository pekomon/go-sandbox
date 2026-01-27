# DungeonDice

DungeonDice is a planned CLI roguelike built around dice-driven combat rounds and class-based runs. The module is scaffolded and ready for tests and implementation.

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

The CLI is under active development; usage examples will land once the run loop is implemented.

---

## Flags

No CLI flags are defined yet.

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
