package stats_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pekomon/go-sandbox/guessr/internal/game"
	"github.com/pekomon/go-sandbox/guessr/internal/stats"
)

func TestStoreLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(stats.EnvPath, filepath.Join(dir, "stats.json"))

	store, err := stats.NewStore()
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got != (game.Stats{}) {
		t.Fatalf("Load() = %+v, want zero stats", got)
	}
}

func TestStoreSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	statsPath := filepath.Join(dir, "stats.json")
	t.Setenv(stats.EnvPath, statsPath)

	store, err := stats.NewStore()
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	want := game.Stats{Games: 3, Wins: 2, TotalGuesses: 9, AverageGuesses: 3}
	if err := store.Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got != want {
		t.Fatalf("round-trip stats = %+v, want %+v", got, want)
	}
}

func TestStoreLoadCorruptedJSON(t *testing.T) {
	dir := t.TempDir()
	statsPath := filepath.Join(dir, "stats.json")
	t.Setenv(stats.EnvPath, statsPath)

	raw := []byte("not-json")
	if err := os.WriteFile(statsPath, raw, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store, err := stats.NewStore()
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}

	if _, err := store.Load(); err == nil {
		t.Fatalf("Load() error = nil, want error")
	}
}

func TestStatsJSONTags(t *testing.T) {
	statsJSON := game.Stats{Games: 1, Wins: 1, TotalGuesses: 5, AverageGuesses: 5.0}
	data, err := json.Marshal(statsJSON)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var roundTrip game.Stats
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if roundTrip != statsJSON {
		t.Fatalf("roundTrip stats = %+v, want %+v", roundTrip, statsJSON)
	}
}
