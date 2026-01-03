package board

import (
	"fmt"
	"math/rand"
)

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

// New creates a fresh board using the provided configuration.
func New(cfg Config) (*Board, error) {
	if cfg.Rows <= 0 || cfg.Cols <= 0 {
		return nil, fmt.Errorf("board: invalid dimensions")
	}
	if cfg.MemeCount < 0 || cfg.MemeCount > cfg.Rows*cfg.Cols {
		return nil, fmt.Errorf("board: invalid meme count")
	}

	cells := make([][]Cell, cfg.Rows)
	for r := 0; r < cfg.Rows; r++ {
		cells[r] = make([]Cell, cfg.Cols)
	}

	placeMemes(cells, cfg)
	setAdjacency(cells)

	return &Board{
		Rows:   cfg.Rows,
		Cols:   cfg.Cols,
		Cells:  cells,
		Status: StatusActive,
	}, nil
}

// Reveal flips the cell at (row, col) and returns the resulting summary.
func (b *Board) Reveal(row, col int) (RevealResult, error) {
	if b == nil {
		return RevealResult{}, fmt.Errorf("board: nil board")
	}
	if row < 0 || row >= b.Rows || col < 0 || col >= b.Cols {
		return RevealResult{}, fmt.Errorf("board: reveal out of bounds")
	}
	if b.Status != StatusActive {
		return RevealResult{Status: b.Status}, nil
	}

	cell := &b.Cells[row][col]
	if cell.Revealed || cell.Flagged {
		return RevealResult{Status: b.Status}, nil
	}

	if cell.HasMeme {
		cell.Revealed = true
		b.Status = StatusLost
		return RevealResult{
			Revealed: []Coord{{Row: row, Col: col}},
			HitMeme:  true,
			Status:   b.Status,
		}, nil
	}

	revealed := floodReveal(b, row, col)
	if isWin(b) {
		b.Status = StatusWon
	}

	return RevealResult{
		Revealed: revealed,
		HitMeme:  false,
		Status:   b.Status,
	}, nil
}

func placeMemes(cells [][]Cell, cfg Config) {
	if cfg.MemeCount == 0 {
		return
	}

	total := cfg.Rows * cfg.Cols
	positions := make([]int, total)
	for i := 0; i < total; i++ {
		positions[i] = i
	}

	rng := rand.New(rand.NewSource(cfg.Seed))
	rng.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < cfg.MemeCount; i++ {
		pos := positions[i]
		r := pos / cfg.Cols
		c := pos % cfg.Cols
		cells[r][c].HasMeme = true
	}
}

func setAdjacency(cells [][]Cell) {
	rows := len(cells)
	if rows == 0 {
		return
	}
	cols := len(cells[0])
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if cells[r][c].HasMeme {
				continue
			}
			cells[r][c].AdjacentMemes = countAdjacentMemes(cells, r, c)
		}
	}
}

func countAdjacentMemes(cells [][]Cell, row, col int) int {
	rows := len(cells)
	cols := len(cells[0])
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			rr := row + dr
			cc := col + dc
			if rr < 0 || rr >= rows || cc < 0 || cc >= cols {
				continue
			}
			if cells[rr][cc].HasMeme {
				count++
			}
		}
	}
	return count
}

func floodReveal(b *Board, row, col int) []Coord {
	type stackItem struct {
		row int
		col int
	}
	stack := []stackItem{{row: row, col: col}}
	revealed := make([]Coord, 0)

	for len(stack) > 0 {
		idx := len(stack) - 1
		item := stack[idx]
		stack = stack[:idx]

		cell := &b.Cells[item.row][item.col]
		if cell.Revealed || cell.Flagged || cell.HasMeme {
			continue
		}
		cell.Revealed = true
		revealed = append(revealed, Coord{Row: item.row, Col: item.col})

		if cell.AdjacentMemes != 0 {
			continue
		}

		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				if dr == 0 && dc == 0 {
					continue
				}
				rr := item.row + dr
				cc := item.col + dc
				if rr < 0 || rr >= b.Rows || cc < 0 || cc >= b.Cols {
					continue
				}
				neighbor := &b.Cells[rr][cc]
				if neighbor.Revealed || neighbor.Flagged || neighbor.HasMeme {
					continue
				}
				stack = append(stack, stackItem{row: rr, col: cc})
			}
		}
	}

	return revealed
}

func isWin(b *Board) bool {
	totalSafe := b.Rows*b.Cols - countMemes(b)
	if totalSafe == 0 {
		return false
	}
	revealedSafe := 0
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			cell := b.Cells[r][c]
			if !cell.HasMeme && cell.Revealed {
				revealedSafe++
			}
		}
	}
	return revealedSafe == totalSafe
}

func countMemes(b *Board) int {
	total := 0
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Cells[r][c].HasMeme {
				total++
			}
		}
	}
	return total
}
