package board

import "fmt"

// Preset defines a difficulty preset for board sizing.
type Preset string

const (
	PresetEasy   Preset = "easy"
	PresetMedium Preset = "medium"
	PresetHard   Preset = "hard"
)

type presetSpec struct {
	rows  int
	cols  int
	memes int
}

var presetSpecs = map[Preset]presetSpec{
	PresetEasy:   {rows: 8, cols: 8, memes: 10},
	PresetMedium: {rows: 16, cols: 16, memes: 40},
	PresetHard:   {rows: 16, cols: 30, memes: 99},
}

// Presets returns the supported difficulty presets in display order.
func Presets() []Preset {
	return []Preset{PresetEasy, PresetMedium, PresetHard}
}

// PresetConfig builds a board Config for the requested preset.
func PresetConfig(preset Preset, seed int64) (Config, error) {
	spec, ok := presetSpecs[preset]
	if !ok {
		return Config{}, fmt.Errorf("board: unknown preset %q", preset)
	}
	return Config{
		Rows:      spec.rows,
		Cols:      spec.cols,
		MemeCount: spec.memes,
		Seed:      seed,
	}, nil
}
