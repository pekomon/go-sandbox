package board_test

import (
	"sort"
	"testing"

	"github.com/pekomon/go-sandbox/memesweeper/internal/board"
)

func TestNewDeterministicLayout(t *testing.T) {
	t.Helper()

	cfg := board.Config{Rows: 4, Cols: 4, MemeCount: 4, Seed: 42}
	b := newBoard(t, cfg)

	gotMemes := collectMemes(b)
	wantMemes := []board.Coord{
		{Row: 1, Col: 3},
		{Row: 2, Col: 3},
		{Row: 3, Col: 0},
		{Row: 3, Col: 3},
	}

	if len(gotMemes) != len(wantMemes) {
		t.Fatalf("board had %d memes; want %d", len(gotMemes), len(wantMemes))
	}

	sortCoords(gotMemes)
	sortCoords(wantMemes)
	for i := range gotMemes {
		if gotMemes[i] != wantMemes[i] {
			t.Fatalf("meme coords mismatch at %d: got %+v want %+v", i, gotMemes[i], wantMemes[i])
		}
	}

	checks := map[board.Coord]int{
		{Row: 0, Col: 2}: 1,
		{Row: 2, Col: 2}: 3,
		{Row: 3, Col: 1}: 1,
	}
	for coord, want := range checks {
		cell := b.Cells[coord.Row][coord.Col]
		if cell.AdjacentMemes != want {
			t.Fatalf("cell %+v AdjacentMemes = %d; want %d", coord, cell.AdjacentMemes, want)
		}
	}
}

func TestRevealWinLossFlow(t *testing.T) {
	t.Helper()

	cfg := board.Config{Rows: 4, Cols: 4, MemeCount: 4, Seed: 42}

	t.Run("hit meme loses", func(t *testing.T) {
		b := newBoard(t, cfg)
		result := reveal(t, b, 3, 3)
		if !result.HitMeme {
			t.Fatalf("Reveal(3,3) should hit a meme")
		}
		if result.Status != board.StatusLost {
			t.Fatalf("status = %v; want %v", result.Status, board.StatusLost)
		}
		if b.Status != board.StatusLost {
			t.Fatalf("board.Status = %v; want %v", b.Status, board.StatusLost)
		}
	})

	t.Run("reveal all safe tiles wins", func(t *testing.T) {
		b := newBoard(t, cfg)
		initial := reveal(t, b, 0, 0)
		if initial.HitMeme {
			t.Fatalf("Reveal(0,0) unexpectedly hit a meme")
		}
		if initial.Status != board.StatusActive {
			t.Fatalf("initial status = %v; want %v", initial.Status, board.StatusActive)
		}
		revealed := map[board.Coord]bool{}
		for _, coord := range initial.Revealed {
			revealed[coord] = true
		}

		safeCells := []board.Coord{
			{Row: 0, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: 2}, {Row: 0, Col: 3},
			{Row: 1, Col: 0}, {Row: 1, Col: 1}, {Row: 1, Col: 2},
			{Row: 2, Col: 0}, {Row: 2, Col: 1}, {Row: 2, Col: 2},
			{Row: 3, Col: 1}, {Row: 3, Col: 2},
		}

		var lastResult board.RevealResult
		for _, coord := range safeCells {
			if revealed[coord] {
				continue
			}
			lastResult = reveal(t, b, coord.Row, coord.Col)
			if lastResult.HitMeme {
				t.Fatalf("Reveal(%d,%d) should not hit a meme", coord.Row, coord.Col)
			}
			for _, newly := range lastResult.Revealed {
				revealed[newly] = true
			}
		}

		if len(revealed) != len(safeCells) {
			t.Fatalf("revealed %d safe cells; want %d", len(revealed), len(safeCells))
		}
		if b.Status != board.StatusWon {
			t.Fatalf("board.Status = %v; want %v", b.Status, board.StatusWon)
		}
		if lastResult.Status != board.StatusWon {
			t.Fatalf("final reveal status = %v; want %v", lastResult.Status, board.StatusWon)
		}
	})
}

func newBoard(t *testing.T, cfg board.Config) *board.Board {
	t.Helper()

	b, err := board.New(cfg)
	if err != nil {
		t.Fatalf("board.New returned error: %v", err)
	}
	return b
}

func reveal(t *testing.T, b *board.Board, row, col int) board.RevealResult {
	t.Helper()

	res, err := b.Reveal(row, col)
	if err != nil {
		t.Fatalf("Reveal(%d,%d) returned error: %v", row, col, err)
	}
	return res
}

func collectMemes(b *board.Board) []board.Coord {
	coords := make([]board.Coord, 0)
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Cells[r][c].HasMeme {
				coords = append(coords, board.Coord{Row: r, Col: c})
			}
		}
	}
	return coords
}

func sortCoords(coords []board.Coord) {
	sort.Slice(coords, func(i, j int) bool {
		if coords[i].Row == coords[j].Row {
			return coords[i].Col < coords[j].Col
		}
		return coords[i].Row < coords[j].Row
	})
}
