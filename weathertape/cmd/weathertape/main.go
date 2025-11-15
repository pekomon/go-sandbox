package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pekomon/go-sandbox/weathertape/internal/forecast"
	"github.com/pekomon/go-sandbox/weathertape/internal/tape"
)

//go:embed sampledata/sample.json
var embeddedSample []byte

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	var (
		sourcePath string
		unitsFlag  string
		width      int
		startStr   string
		endStr     string
	)
	fs := flag.NewFlagSet("weathertape", flag.ContinueOnError)
	fs.StringVar(&sourcePath, "source", "", "path to the forecast JSON file")
	fs.StringVar(&unitsFlag, "units", "metric", "temperature units: metric or imperial")
	fs.IntVar(&width, "width", 10, "width of the temperature bar graph")
	fs.StringVar(&startStr, "start", "", "optional RFC3339 start time filter")
	fs.StringVar(&endStr, "end", "", "optional RFC3339 end time filter")
	fs.SetOutput(new(nopWriter))
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "invalid flags")
		return 2
	}

	units, err := parseUnits(unitsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	if width < 5 {
		fmt.Fprintln(os.Stderr, "width must be >=5")
		return 2
	}

	startTime, endTime, err := parseRange(startStr, endStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if sourcePath == "" {
		sourcePath = os.Getenv("WEATHERTAPE_DATA")
	}

	var entries []forecast.Entry
	if strings.TrimSpace(sourcePath) != "" {
		path := expandHome(strings.TrimSpace(sourcePath))
		entries, err = forecast.LoadFile(path)
	} else {
		entries, err = forecast.LoadBytes(embeddedSample, "embedded:sample.json")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	entries = applyRange(entries, startTime, endTime)
	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "no forecast rows matched the requested range")
		return 1
	}

	rendered, err := tape.Render(entries, tape.Options{
		Units: units,
		Width: width,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Fprint(os.Stdout, rendered)
	return 0
}

func parseUnits(flagValue string) (tape.Units, error) {
	switch strings.ToLower(strings.TrimSpace(flagValue)) {
	case "", "metric", "c", "celsius":
		return tape.UnitsMetric, nil
	case "imperial", "f", "fahrenheit":
		return tape.UnitsImperial, nil
	default:
		return tape.UnitsMetric, fmt.Errorf("invalid units %q (use metric or imperial)", flagValue)
	}
}

func parseRange(startStr, endStr string) (*time.Time, *time.Time, error) {
	var (
		startPtr *time.Time
		endPtr   *time.Time
	)
	if strings.TrimSpace(startStr) != "" {
		ts, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid start time %q (use RFC3339): %w", startStr, err)
		}
		startPtr = &ts
	}
	if strings.TrimSpace(endStr) != "" {
		ts, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid end time %q (use RFC3339): %w", endStr, err)
		}
		endPtr = &ts
	}
	if startPtr != nil && endPtr != nil && startPtr.After(*endPtr) {
		return nil, nil, fmt.Errorf("start time must be <= end time")
	}
	return startPtr, endPtr, nil
}

func applyRange(entries []forecast.Entry, start, end *time.Time) []forecast.Entry {
	if start == nil && end == nil {
		return entries
	}
	filtered := make([]forecast.Entry, 0, len(entries))
	for _, entry := range entries {
		if start != nil && entry.Time.Before(*start) {
			continue
		}
		if end != nil && entry.Time.After(*end) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return home + path[1:]
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }
