package forecast_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/pekomon/go-sandbox/weathertape/internal/forecast"
)

func TestLoadFileSuccess(t *testing.T) {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", "sample.json")
	entries, err := forecast.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile(%q) returned error: %v", path, err)
	}
	if len(entries) != 4 {
		t.Fatalf("LoadFile(%q) returned %d entries; want 4", path, len(entries))
	}

	first := entries[0]
	wantTime := time.Date(2025, 2, 10, 9, 0, 0, 0, time.UTC)
	if !first.Time.Equal(wantTime) {
		t.Fatalf("first entry time = %s; want %s", first.Time.Format(time.RFC3339), wantTime.Format(time.RFC3339))
	}
	if first.TempC != 12.0 {
		t.Fatalf("first entry TempC = %.1f; want 12.0", first.TempC)
	}
	if first.PrecipPercent != 10 {
		t.Fatalf("first entry PrecipPercent = %d; want 10", first.PrecipPercent)
	}
	if first.WindKPH != 8.0 {
		t.Fatalf("first entry WindKPH = %.1f; want 8.0", first.WindKPH)
	}
	if first.WindDirection != "NE" {
		t.Fatalf("first entry WindDirection = %q; want %q", first.WindDirection, "NE")
	}

	last := entries[len(entries)-1]
	if last.TempC != 21.0 {
		t.Fatalf("last entry TempC = %.1f; want 21.0", last.TempC)
	}
}

func TestLoadFileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	if _, err := forecast.LoadFile(path); err == nil {
		t.Fatalf("LoadFile(%q) returned nil error; want failure for missing file", path)
	}
}
