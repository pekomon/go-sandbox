package game

import (
	"errors"
	"math/rand"
)

type Phase int

const (
	PhaseAttract Phase = iota
	PhaseShowing
	PhaseInput
	PhaseGameOver
)

type Input struct {
	Start   bool
	Restart bool
	Pad     int
	Press   bool
}

type Config struct {
	Width             float64
	Height            float64
	PadCount          int
	StartLength       int
	ShowTicks         int
	GapTicks          int
	InputTimeoutTicks int
	Seed              int64
}

type State struct {
	Width             float64
	Height            float64
	PadCount          int
	Phase             Phase
	Round             int
	Score             int
	Sequence          []int
	ShowIndex         int
	InputIndex        int
	LitPad            int
	ShowTicks         int
	GapTicks          int
	InputTimeoutTicks int
	PhaseTick         int
	WatchdogTick      int

	baseShowTicks         int
	baseGapTicks          int
	baseInputTimeoutTicks int
	startLength           int
	rng                   *rand.Rand
}

func DefaultConfig() Config {
	return Config{
		Width:             640,
		Height:            360,
		PadCount:          4,
		StartLength:       2,
		ShowTicks:         24,
		GapTicks:          14,
		InputTimeoutTicks: 180,
		Seed:              1,
	}
}

func New(cfg Config) (*State, error) {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, errors.New("invalid dimensions")
	}
	if cfg.PadCount < 2 {
		return nil, errors.New("invalid pad count")
	}
	if cfg.StartLength < 1 {
		return nil, errors.New("invalid start length")
	}
	if cfg.ShowTicks < 1 || cfg.GapTicks < 1 || cfg.InputTimeoutTicks < 1 {
		return nil, errors.New("invalid timing configuration")
	}

	return &State{
		Width:                 cfg.Width,
		Height:                cfg.Height,
		PadCount:              cfg.PadCount,
		Phase:                 PhaseAttract,
		LitPad:                -1,
		ShowTicks:             cfg.ShowTicks,
		GapTicks:              cfg.GapTicks,
		InputTimeoutTicks:     cfg.InputTimeoutTicks,
		baseShowTicks:         cfg.ShowTicks,
		baseGapTicks:          cfg.GapTicks,
		baseInputTimeoutTicks: cfg.InputTimeoutTicks,
		startLength:           cfg.StartLength,
		rng:                   rand.New(rand.NewSource(cfg.Seed)),
	}, nil
}

func (s *State) Step(Input) error {
	return nil
}
