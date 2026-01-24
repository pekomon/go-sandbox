# ThumbForge

ThumbForge is a CLI for batch thumbnail generation. It resizes PNG/JPEG inputs to fixed-size thumbnails and writes PNG or JPEG outputs offline.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd thumbforge
make deps   # optional; runs go mod tidy
make build  # produces ./bin/thumbforge
```

If you prefer raw Go commands:

```bash
cd thumbforge
go build -o bin/thumbforge ./cmd/thumbforge
```

---

## Usage

Generate thumbnails using a WxH size:

```bash
./bin/thumbforge --in ./photos --out ./thumbs --size 320x240
```

Generate thumbnails using width/height flags:

```bash
./bin/thumbforge --in ./photos --out ./thumbs --width 320 --height 240 --format jpg
```

### Flags

| Flag | Description | Default |
| ---- | ----------- | ------- |
| `--in` | Input directory (required). | _none_ |
| `--out` | Output directory (required). | _none_ |
| `--size` | Thumbnail size in `WxH` form. | _none_ |
| `--width` | Thumbnail width in pixels (use with `--height`). | `0` |
| `--height` | Thumbnail height in pixels (use with `--width`). | `0` |
| `--format` | Output format (`png`, `jpg`, `jpeg`). | `png` |
| `--crop` | Accepted for future use (currently no effect). | `false` |

Notes:
- Provide either `--size` or `--width` + `--height` (not both).
- Output files keep the input base name with the output format extension.
- Empty input directories return an error.

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
- The CLI lives in `cmd/thumbforge` and implementation packages under `internal/`.
- Follow the repository conventions in [../agents.md](../agents.md) for branching strategy, PR templates, and release cadence.
