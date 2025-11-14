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

The CLI entry point will live under `cmd/weathertape`. Once the feature PRs land you will be able to run the tool via:

```bash
./bin/weathertape --source ./testdata/sample.json
```

Environment overrides (subject to change during implementation):

- `WEATHERTAPE_DATA` — absolute/relative path to the forecast data file
- `WEATHERTAPE_CACHE` — override the cache directory (default `~/.weathertape`)

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
