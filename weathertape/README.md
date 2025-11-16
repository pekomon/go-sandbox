# WeatherTape

WeatherTape is a terminal dashboard that turns structured hourly forecast data into an ASCII tape. It loads JSON forecasts (either embedded sample data or a user-provided file), filters the hours you care about, and renders aligned rows that show the hour, temperature, precipitation chance, wind speed/direction, and a proportional temperature bar.

---

## Installation

Requirements: Go 1.25+, `make`, and a POSIX shell.

```bash
# From repo root
cd weathertape
make deps   # optional; runs go mod tidy
make build  # produces ./bin/weathertape
```

If you prefer raw Go commands:

```bash
cd weathertape
go build -o bin/weathertape ./cmd/weathertape
```

---

## Usage

With no flags, WeatherTape renders the embedded `cmd/weathertape/sampledata/sample.json` payload:

```bash
cd weathertape
./bin/weathertape
```

Typical invocations:

```bash
# Use a custom dataset with imperial units and a wider bar graph
./bin/weathertape --source ./testdata/hourly.json --units imperial --width 20

# Narrow the tape to a specific RFC3339 time window
./bin/weathertape \
  --start 2025-02-10T10:00:00Z \
  --end   2025-02-10T14:00:00Z
```

Sample output (metric, width=10):

```
Hour  Temp  Trend       Precip Wind
----  ----  ----------  ------ ----
09:00 12°C  █░░░░░░░░░░   10%  NE8kph
10:00 15°C  ███░░░░░░░░   40%  E11kph
11:00 19°C  ███████░░░░   70%  SE20kph
12:00 21°C  ██████████   15%  S25kph
```

### CLI flags

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `--source` | *(empty)* | Optional path to a JSON forecast file. When omitted, WeatherTape uses `WEATHERTAPE_DATA` (if set) and finally the embedded sample. |
| `--units` | `metric` | Temperature/wind units (`metric`, `imperial`, `c`, `f`, etc.). |
| `--width` | `10` | Width (characters) of the temperature bar. Minimum accepted value is 5. |
| `--start` | *(empty)* | RFC3339 timestamp that acts as the inclusive start of the rendered range. |
| `--end` | *(empty)* | RFC3339 timestamp that acts as the inclusive end of the rendered range. |

### Data format

WeatherTape expects a JSON array whose entries match:

```json
{
  "hour": "2025-02-10T09:00:00Z",
  "temp_c": 12.0,
  "precip_pct": 10,
  "wind_kph": 8.0,
  "wind_dir": "NE"
}
```

The loader validates the timestamp (`hour`) and rejects empty datasets.

---

## Environment variables

| Variable | Purpose |
| -------- | ------- |
| `WEATHERTAPE_DATA` | Absolute or relative path to a forecast file that becomes the default when `--source` is not provided. |

---

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0` | Success: forecast rows rendered to stdout. |
| `1` | Runtime failure (file I/O, parsing errors, empty result set after filtering). |
| `2` | Invalid CLI usage (flag parsing errors, unsupported units, invalid `--start/--end`, or `--width < 5`). |

Normal output is printed to stdout; all error messages go to stderr.

---

## Testing & coverage

```bash
cd weathertape
make test   # go test ./...
make cover  # go test ./... -coverprofile=cover.out && go tool cover -func cover.out
```

The GitHub Actions workflow mirrors these targets and uploads the `cover.out` artifact for pull requests touching this module.

---

## Development notes

- Modules stick to the Go standard library; `internal/tape` owns rendering and `internal/forecast` owns loading/validation.
- The CLI lives in `cmd/weathertape` and embeds a small dataset for demo/testing.
- Follow the repository conventions in [../agents.md](../agents.md) for branching strategy, PR templates, and release cadence.
