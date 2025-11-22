package board

import "errors"

// Status represents the state of a board as the game progresses.
type Status int

const (
	// StatusActive indicates the game is still running.
	StatusActive Status = iota
	// StatusWon indicates every safe cell has been revealed.
	StatusWon
	// StatusLost indicates a meme was revealed.
	StatusLost
)

// Config controls the size and entropy used when generating a new board.
type Config struct {
	Rows      int
	Cols      int
	MemeCount int
	Seed      int64
}

// Coord identifies a row/column cell location.
type Coord struct {
	Row int
	Col int
}

// Cell describes a tile on the board.
type Cell struct {
	HasMeme       bool
	AdjacentMemes int
	Revealed      bool
	Flagged       bool
}

// Board models the full playfield state.
type Board struct {
	Rows   int
	Cols   int
	Cells  [][]Cell
	Status Status
}

// RevealResult summarizes what happened after revealing a cell.
type RevealResult struct {
	Revealed []Coord
	HitMeme  bool
	Status   Status
}

// ErrNotImplemented indicates the board engine hasn't been built yet.
var ErrNotImplemented = errors.New("board: not implemented")

// New creates a fresh board using the provided configuration.
func New(cfg Config) (*Board, error) {
	return nil, ErrNotImplemented
}

// Reveal flips the cell at (row, col) and returns the resulting summary.
func (b *Board) Reveal(row, col int) (RevealResult, error) {
	return RevealResult{}, ErrNotImplemented
}
