package board_test

import (
	"testing"

	"github.com/pekomon/go-sandbox/memesweeper/internal/board"
)

func TestPresetConfigValues(t *testing.T) {
	t.Helper()

	cases := []struct {
		name      string
		preset    board.Preset
		rows      int
		cols      int
		memeCount int
	}{
		{name: "easy", preset: board.PresetEasy, rows: 8, cols: 8, memeCount: 10},
		{name: "medium", preset: board.PresetMedium, rows: 16, cols: 16, memeCount: 40},
		{name: "hard", preset: board.PresetHard, rows: 16, cols: 30, memeCount: 99},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := board.PresetConfig(tc.preset, 99)
			if err != nil {
				t.Fatalf("PresetConfig(%q) error: %v", tc.preset, err)
			}
			if cfg.Rows != tc.rows || cfg.Cols != tc.cols || cfg.MemeCount != tc.memeCount {
				t.Fatalf("PresetConfig(%q) = %+v; want rows=%d cols=%d memes=%d", tc.preset, cfg, tc.rows, tc.cols, tc.memeCount)
			}
			if cfg.Seed != 99 {
				t.Fatalf("PresetConfig(%q) seed=%d; want 99", tc.preset, cfg.Seed)
			}
		})
	}
}

func TestPresetConfigDeterministic(t *testing.T) {
	t.Helper()

	for _, preset := range board.Presets() {
		cfg, err := board.PresetConfig(preset, 42)
		if err != nil {
			t.Fatalf("PresetConfig(%q) error: %v", preset, err)
		}
		first := newBoard(t, cfg)
		second := newBoard(t, cfg)

		got := collectMemes(first)
		want := collectMemes(second)
		sortCoords(got)
		sortCoords(want)

		if len(got) != len(want) {
			t.Fatalf("%q meme count mismatch: got %d want %d", preset, len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("%q meme coords mismatch at %d: got %+v want %+v", preset, i, got[i], want[i])
			}
		}
	}
}

func TestPresetConfigUnknown(t *testing.T) {
	t.Helper()

	if _, err := board.PresetConfig(board.Preset("nope"), 0); err == nil {
		t.Fatalf("PresetConfig should fail for unknown preset")
	}
}
