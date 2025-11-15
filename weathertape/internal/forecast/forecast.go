package forecast

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry models a single hourly forecast row used by the renderer.
type Entry struct {
	Time           time.Time
	TempC          float64
	PrecipPercent  int
	WindKPH        float64
	WindDirection  string
	SourceFilePath string
}

type rawEntry struct {
	Hour          string  `json:"hour"`
	TempC         float64 `json:"temp_c"`
	PrecipPercent int     `json:"precip_pct"`
	WindKPH       float64 `json:"wind_kph"`
	WindDirection string  `json:"wind_dir"`
}

// LoadFile loads forecast data from the given file path.
func LoadFile(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadBytes(data, path)
}

// LoadBytes parses forecast data from the provided JSON payload.
func LoadBytes(data []byte, source string) ([]Entry, error) {
	var raw []rawEntry
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("forecast: decode %s: %w", source, err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("forecast: no entries in %s", source)
	}

	entries := make([]Entry, len(raw))
	for i, r := range raw {
		ts, err := time.Parse(time.RFC3339, r.Hour)
		if err != nil {
			return nil, fmt.Errorf("forecast: parse hour %q (row %d): %w", r.Hour, i, err)
		}
		entries[i] = Entry{
			Time:           ts,
			TempC:          r.TempC,
			PrecipPercent:  r.PrecipPercent,
			WindKPH:        r.WindKPH,
			WindDirection:  r.WindDirection,
			SourceFilePath: source,
		}
	}
	return entries, nil
}
