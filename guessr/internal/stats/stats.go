package stats

import "github.com/pekomon/go-sandbox/guessr/internal/game"

// Path resolution is env-overridable for tests:
//
//	GUESSR_STATS_PATH=/tmp/whatever.json
const EnvPath = "GUESSR_STATS_PATH"

// Store is a file-backed implementation of game.StatsStore.
// STUB for now: methods return ErrNotImplemented.
type Store struct {
	// jsonPath string
}

func NewStore() (*Store, error) { return &Store{}, nil }

func (s *Store) Load() (game.Stats, error) { return game.Stats{}, game.ErrNotImplemented }
func (s *Store) Save(st game.Stats) error  { return game.ErrNotImplemented }
