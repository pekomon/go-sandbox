package tape

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

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

// Render builds an ASCII tape using the provided forecast entries.
func Render(entries []forecast.Entry, opts Options) (string, error) {
	width := opts.Width
	if width <= 0 {
		width = 10
	}

	var b strings.Builder
	effectiveWidth := width
	if effectiveWidth < 10 {
		effectiveWidth = 10
	}
	trendHeader := padRight("Trend", effectiveWidth)
	trendRule := strings.Repeat("-", effectiveWidth)
	fmt.Fprintf(&b, "Hour  Temp  %s  Precip Wind\n", trendHeader)
	fmt.Fprintf(&b, "----  ----  %s  ------ ----\n", trendRule)

	if len(entries) == 0 {
		return b.String(), nil
	}

	minTemp, maxTemp := minMaxTemps(entries)
	for _, entry := range entries {
		bar := buildBar(entry.TempC, minTemp, maxTemp, width)
		hour := entry.Time.Format("15:04")
		tempStr := padRight(formatTemp(entry.TempC, opts.Units), 4)
		precip := fmt.Sprintf("%d%%", entry.PrecipPercent)
		wind := formatWind(entry, opts.Units)
		fmt.Fprintf(&b, "%s %s  %s   %s  %s\n", hour, tempStr, bar, precip, wind)
	}

	return b.String(), nil
}

func minMaxTemps(entries []forecast.Entry) (float64, float64) {
	min := entries[0].TempC
	max := entries[0].TempC
	for _, e := range entries[1:] {
		if e.TempC < min {
			min = e.TempC
		}
		if e.TempC > max {
			max = e.TempC
		}
	}
	return min, max
}

func buildBar(temp, min, max float64, width int) string {
	if width <= 0 {
		width = 10
	}
	if max == min {
		return strings.Repeat("█", width)
	}
	ratio := (temp - min) / (max - min)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(math.Round(ratio * float64(width)))
	if filled < 1 {
		filled = 1
	}
	if filled > width {
		filled = width
	}
	if filled == width {
		return strings.Repeat("█", width)
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width+1-filled)
}

func formatTemp(tempC float64, units Units) string {
	switch units {
	case UnitsImperial:
		tempF := math.Round(tempC*9/5 + 32)
		return fmt.Sprintf("%d°F", int(tempF))
	default:
		return fmt.Sprintf("%d°C", int(math.Round(tempC)))
	}
}

func formatWind(entry forecast.Entry, units Units) string {
	speed := entry.WindKPH
	unit := "kph"
	if units == UnitsImperial {
		speed = entry.WindKPH * 0.621371
		unit = "mph"
	}
	speed = math.Round(speed)
	return fmt.Sprintf("%s%d%s", entry.WindDirection, int(speed), unit)
}

func padRight(s string, width int) string {
	if width <= 0 {
		return s
	}
	runes := utf8.RuneCountInString(s)
	if runes >= width {
		return s
	}
	return s + strings.Repeat(" ", width-runes)
}
