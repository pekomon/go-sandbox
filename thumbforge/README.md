# ThumbForge

ThumbForge is a CLI for batch thumbnail generation. It will resize/crop images, preserve EXIF-safe metadata, and export fixed-size assets offline. This module is currently scaffolded; the CLI and tests will land in upcoming issues.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd thumbforge
make deps   # optional; runs go mod tidy
make build  # produces ./bin/thumbforge once the CLI exists
```

If you prefer raw Go commands:

```bash
cd thumbforge
go build -o bin/thumbforge ./cmd/thumbforge
```

---

## Usage

ThumbForge is not implemented yet. Planned usage will look like:

```bash
./bin/thumbforge --in ./photos --out ./thumbs --size 320x240
```

---

## Environment variables

No environment variables are defined yet.

---

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Success. |
| `1` | Runtime failure (I/O or processing errors). |
| `2` | Invalid CLI usage. |

Normal output will be printed to stdout; all error messages go to stderr.

---

## Testing & coverage

```bash
cd thumbforge
make test   # go test ./...
make cover  # go test ./... -coverprofile=cover.out && go tool cover -func cover.out
```

The GitHub Actions workflow mirrors these targets and uploads the `cover.out` artifact for pull requests touching this module.

---

## Development notes

- Modules stick to the Go standard library; image processing dependencies will be justified if needed.
- The CLI will live in `cmd/thumbforge` and implementation packages under `internal/`.
- Follow the repository conventions in [../agents.md](../agents.md) for branching strategy, PR templates, and release cadence.
