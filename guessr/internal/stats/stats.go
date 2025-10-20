package stats

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/pekomon/go-sandbox/guessr/internal/game"
)

const (
	EnvPath        = "GUESSR_STATS_PATH"
	defaultRelPath = ".guessr/stats.json"
)

// Store persists game.Stats to a JSON file.
type Store struct {
	jsonPath string
}

func NewStore() (*Store, error) {
	if p := os.Getenv(EnvPath); p != "" {
		return &Store{jsonPath: p}, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{jsonPath: filepath.Join(home, defaultRelPath)}, nil
}

func (s *Store) path() string { return s.jsonPath }

func (s *Store) Load() (game.Stats, error) {
	b, err := os.ReadFile(s.path())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return game.Stats{}, nil // not found => zero stats
		}
		return game.Stats{}, err
	}
	var st game.Stats
	if err := json.Unmarshal(b, &st); err != nil {
		return game.Stats{}, err
	}
	return st, nil
}

func (s *Store) Save(st game.Stats) error {
	dir := filepath.Dir(s.path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// ensure AverageGuesses coherent
	if st.Games > 0 && st.AverageGuesses == 0 && st.TotalGuesses > 0 {
		st.AverageGuesses = float64(st.TotalGuesses) / float64(st.Games)
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path())
}
