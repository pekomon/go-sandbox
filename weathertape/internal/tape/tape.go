package tape

import (
	"errors"

	"github.com/pekomon/go-sandbox/weathertape/internal/forecast"
)

// Units controls whether the tape renders metric or imperial measurements.
type Units int

const (
	// UnitsMetric renders °C temperatures and kph wind speeds.
	UnitsMetric Units = iota
	// UnitsImperial renders °F temperatures and mph wind speeds.
	UnitsImperial
)

// Options configures how the ASCII tape should be rendered.
type Options struct {
	Units Units
	Width int // number of characters used for the bar graph section
}

// ErrNotImplemented signals the renderer still needs an implementation.
var ErrNotImplemented = errors.New("tape: renderer not implemented")

// Render builds an ASCII tape using the provided forecast entries.
func Render(entries []forecast.Entry, opts Options) (string, error) {
	return "", ErrNotImplemented
}
