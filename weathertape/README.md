# WeatherTape

WeatherTape is a terminal dashboard that renders hourly forecast data as an ASCII "tape" with temperature bars and precipitation markers. The module is scaffolded here; implementation will follow the repo's tests-first workflow.

---

## Planned Features

- Load deterministic weather samples from JSON/CSV files (offline only).
- Render a scrolling ASCII tape of hours, temperatures, precipitation chance, and wind arrows.
- Provide CLI flags for timeframe selection and units (metric/imperial).
- Persist recent runs to `~/.weathertape/` for later comparison (configurable via env var).

---

## Installation

```bash
# From repo root
cd weathertape
make build
# Binary will be created at ./bin/weathertape
```

Or directly with Go:

```bash
cd weathertape
go build -o bin/weathertape ./cmd/weathertape
```

## Usage

```bash
# Read a custom data file
./bin/weathertape --source ./fixtures/hourly.json --units imperial

# Use RFC3339 range filters
./bin/weathertape --start 2025-02-10T10:00:00Z --end 2025-02-10T12:00:00Z
```

Flags:

| Flag | Description |
| ---- | ----------- |
| `--source` | Path to a JSON forecast file. Defaults to `WEATHERTAPE_DATA`. When unset, the binary falls back to the embedded sample dataset (taken from `cmd/weathertape/sampledata/sample.json`). |
| `--units` | Either `metric` (°C / kph, default) or `imperial` (°F / mph). |
| `--width` | Width (in characters) of the ASCII temperature bar. Minimum 5, default 10. |
| `--start`, `--end` | Optional RFC3339 timestamps to filter the hourly rows. |

Environment overrides:

- `WEATHERTAPE_DATA` — Absolute/relative path to a data file (used when `--source` is omitted).
- `WEATHERTAPE_CACHE` — Override the cache directory (default `~/.weathertape`); reserved for future persistence work.

## Testing & Coverage

```bash
cd weathertape
make test
make cover
```

These targets will execute the module's unit/integration tests once they are added. CI runs the same commands and publishes coverage artifacts for PRs touching `weathertape/`.

## Development Notes

- Go 1.25 baseline; no external dependencies planned.
- Logic will live under `internal/` packages, leaving `cmd/weathertape` as glue code.
- Follow the shared conventions in [../agents.md](../agents.md) for branching, tests-first flow, and CI additions.
