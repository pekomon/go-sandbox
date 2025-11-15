package forecast

import (
	"errors"
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

// ErrNotImplemented indicates the loader has not been implemented yet.
var ErrNotImplemented = errors.New("forecast: loader not implemented")

// LoadFile loads forecast data from a JSON/CSV file path and returns normalized entries.
func LoadFile(path string) ([]Entry, error) {
	return nil, ErrNotImplemented
}
