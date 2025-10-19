package game

import (
	"io"
)

type Options struct {
	Max      int // default: 100
	Attempts int // default: 7
	Seed     int64
}

var ErrNotImplemented = errNotImplemented("not implemented")

type errNotImplemented string

func (e errNotImplemented) Error() string { return string(e) }

// Run plays one game by reading guesses (one per line) from r and writing hints/results to w.
// It should respect Options and use the provided random seed (deterministic).
// STUB for now: return ErrNotImplemented so tests fail.
func Run(r io.Reader, w io.Writer, opts Options, stats StatsStore) error {
	return ErrNotImplemented
}

// StatsStore abstracts persistence (implemented in internal/stats).
type StatsStore interface {
	Load() (Stats, error)
	Save(Stats) error
}

// Stats describes game statistics.
type Stats struct {
	Games          int     `json:"games"`
	Wins           int     `json:"wins"`
	TotalGuesses   int     `json:"total_guesses"`
	AverageGuesses float64 `json:"average_guesses"`
}
