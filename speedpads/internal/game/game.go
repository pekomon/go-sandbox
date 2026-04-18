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
	BufferedPad       int

	baseShowTicks         int
	baseGapTicks          int
	baseInputTimeoutTicks int
	startLength           int
	rng                   *rand.Rand
	bufferedPress         bool
	pendingRoundAdvance   bool
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
		BufferedPad:           -1,
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

func (s *State) Step(in Input) error {
	if s == nil {
		return errors.New("nil state")
	}

	switch s.Phase {
	case PhaseAttract:
		if s.startRequested(in) {
			return s.beginRound(true)
		}
	case PhaseShowing:
		s.advanceShowPhase()
		if s.Phase == PhaseInput {
			if in.Press {
				s.advanceInputPhase(in)
			}
		}
	case PhaseInput:
		s.advanceInputPhase(in)
	case PhaseGameOver:
		if s.restartRequested(in) {
			s.reset()
		}
	}

	return nil
}

func (s *State) startRequested(in Input) bool {
	return in.Start || in.Restart
}

func (s *State) restartRequested(in Input) bool {
	return in.Restart || in.Start
}

func (s *State) reset() {
	s.Phase = PhaseAttract
	s.Round = 0
	s.Score = 0
	s.Sequence = s.Sequence[:0]
	s.ShowIndex = 0
	s.InputIndex = 0
	s.LitPad = -1
	s.BufferedPad = -1
	s.ShowTicks = s.baseShowTicks
	s.GapTicks = s.baseGapTicks
	s.InputTimeoutTicks = s.baseInputTimeoutTicks
	s.PhaseTick = 0
	s.WatchdogTick = 0
	s.bufferedPress = false
	s.pendingRoundAdvance = false
}

func (s *State) beginRound(fresh bool) error {
	if fresh {
		s.reset()
		s.Round = 1
		for i := 0; i < s.startLength; i++ {
			s.Sequence = append(s.Sequence, s.nextPad())
		}
	} else {
		s.Round++
		s.Sequence = append(s.Sequence, s.nextPad())
		s.applyDifficultyRamp()
	}

	s.Phase = PhaseShowing
	s.ShowIndex = 0
	s.InputIndex = 0
	s.PhaseTick = 0
	s.WatchdogTick = 0
	s.bufferedPress = false
	s.pendingRoundAdvance = false
	s.BufferedPad = -1
	s.LitPad = s.Sequence[0]
	return nil
}

func (s *State) nextPad() int {
	return s.rng.Intn(s.PadCount)
}

func (s *State) applyDifficultyRamp() {
	showTicks := s.baseShowTicks - (s.Round - 1)
	if showTicks < 1 {
		showTicks = 1
	}

	gapTicks := s.baseGapTicks - ((s.Round - 1) / 2) - ((s.Round - 1) % 2)
	if gapTicks < 0 {
		gapTicks = 0
	}

	timeoutTicks := s.baseInputTimeoutTicks - ((s.Round - 1) * 2)
	if timeoutTicks < 1 {
		timeoutTicks = 1
	}

	s.ShowTicks = showTicks
	s.GapTicks = gapTicks
	s.InputTimeoutTicks = timeoutTicks
}

func (s *State) advanceShowPhase() {
	s.PhaseTick++

	if s.LitPad >= 0 {
		if s.PhaseTick < s.ShowTicks {
			return
		}
		s.LitPad = -1
		s.PhaseTick = 0
		return
	}

	if s.PhaseTick < s.GapTicks {
		return
	}

	s.ShowIndex++
	s.PhaseTick = 0
	if s.ShowIndex >= len(s.Sequence) {
		s.Phase = PhaseInput
		s.InputIndex = 0
		s.WatchdogTick = 0
		s.LitPad = -1
		return
	}

	s.LitPad = s.Sequence[s.ShowIndex]
}

func (s *State) advanceInputPhase(in Input) {
	if s.pendingRoundAdvance {
		if s.PhaseTick > 0 {
			s.PhaseTick--
			if s.PhaseTick == 0 {
				s.LitPad = -1
				s.pendingRoundAdvance = false
				_ = s.beginRound(false)
			}
		}
		return
	}

	if in.Press {
		if in.Pad < 0 || in.Pad >= s.PadCount {
			s.Phase = PhaseGameOver
			s.LitPad = -1
			return
		}

		s.LitPad = in.Pad
		s.PhaseTick = 1
		s.WatchdogTick = 0
		if in.Pad != s.Sequence[s.InputIndex] {
			s.Phase = PhaseGameOver
			return
		}

		s.Score++
		s.InputIndex++
		if s.InputIndex >= len(s.Sequence) {
			s.pendingRoundAdvance = true
		}
		return
	}

	s.WatchdogTick++
	if s.WatchdogTick >= s.InputTimeoutTicks {
		s.Phase = PhaseGameOver
		s.LitPad = -1
		return
	}

	if s.PhaseTick > 0 {
		s.PhaseTick--
		if s.PhaseTick == 0 {
			s.LitPad = -1
		}
	}
}
