package tape_test

import (
	"strings"
	"testing"
	"time"

	"github.com/pekomon/go-sandbox/weathertape/internal/forecast"
	"github.com/pekomon/go-sandbox/weathertape/internal/tape"
)

func TestRenderMetricTape(t *testing.T) {
	t.Helper()

	entries := sampleEntries()
	opts := tape.Options{Units: tape.UnitsMetric, Width: 10}
	got, err := tape.Render(entries, opts)
	if err != nil {
		t.Fatalf("Render(metric) returned error: %v", err)
	}

	const want = "" +
		"Hour  Temp  Trend       Precip Wind\n" +
		"----  ----  ----------  ------ ----\n" +
		"09:00 12°C  █░░░░░░░░░░   10%  NE8kph\n" +
		"10:00 15°C  ███░░░░░░░░   40%  E11kph\n" +
		"11:00 19°C  ███████░░░░   70%  SE20kph\n" +
		"12:00 21°C  ██████████   15%  S25kph"

	if strings.TrimSpace(got) != want {
		t.Fatalf("Render(metric) mismatch.\nGot:\n%s\n\nWant:\n%s", got, want)
	}
}

func TestRenderImperialTape(t *testing.T) {
	t.Helper()

	entries := sampleEntries()
	opts := tape.Options{Units: tape.UnitsImperial, Width: 10}
	got, err := tape.Render(entries, opts)
	if err != nil {
		t.Fatalf("Render(imp) returned error: %v", err)
	}

	const want = "" +
		"Hour  Temp  Trend       Precip Wind\n" +
		"----  ----  ----------  ------ ----\n" +
		"09:00 54°F  █░░░░░░░░░░   10%  NE5mph\n" +
		"10:00 59°F  ███░░░░░░░░   40%  E7mph\n" +
		"11:00 65°F  ███████░░░░   70%  SE12mph\n" +
		"12:00 70°F  ██████████   15%  S16mph"

	if strings.TrimSpace(got) != want {
		t.Fatalf("Render(imp) mismatch.\nGot:\n%s\n\nWant:\n%s", got, want)
	}
}

func sampleEntries() []forecast.Entry {
	return []forecast.Entry{
		{
			Time:          time.Date(2025, 2, 10, 9, 0, 0, 0, time.UTC),
			TempC:         12.0,
			PrecipPercent: 10,
			WindKPH:       8,
			WindDirection: "NE",
		},
		{
			Time:          time.Date(2025, 2, 10, 10, 0, 0, 0, time.UTC),
			TempC:         15.0,
			PrecipPercent: 40,
			WindKPH:       11,
			WindDirection: "E",
		},
		{
			Time:          time.Date(2025, 2, 10, 11, 0, 0, 0, time.UTC),
			TempC:         18.5,
			PrecipPercent: 70,
			WindKPH:       20,
			WindDirection: "SE",
		},
		{
			Time:          time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC),
			TempC:         21.0,
			PrecipPercent: 15,
			WindKPH:       25,
			WindDirection: "S",
		},
	}
}
