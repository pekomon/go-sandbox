package game_test

import (
	"testing"

	"github.com/pekomon/go-sandbox/speedpads/internal/game"
)

func newConfig() game.Config {
	cfg := game.DefaultConfig()
	cfg.ShowTicks = 2
	cfg.GapTicks = 1
	cfg.InputTimeoutTicks = 5
	cfg.Seed = 7
	return cfg
}

func newState(t *testing.T) *game.State {
	t.Helper()

	st, err := game.New(newConfig())
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	return st
}

func stepMany(t *testing.T, st *game.State, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := st.Step(game.Input{}); err != nil {
			t.Fatalf("Step %d: %v", i, err)
		}
	}
}

func TestNewStartsInAttractPhase(t *testing.T) {
	t.Parallel()

	st := newState(t)

	if st.Phase != game.PhaseAttract {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseAttract)
	}
	if st.Round != 0 {
		t.Fatalf("round = %d, want 0", st.Round)
	}
	if st.LitPad != -1 {
		t.Fatalf("lit pad = %d, want -1", st.LitPad)
	}
	if len(st.Sequence) != 0 {
		t.Fatalf("sequence length = %d, want 0", len(st.Sequence))
	}
}

func TestStartBeginsShowingDeterministicSequence(t *testing.T) {
	t.Parallel()

	st := newState(t)
	if err := st.Step(game.Input{Start: true}); err != nil {
		t.Fatalf("Step start: %v", err)
	}

	if st.Phase != game.PhaseShowing {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseShowing)
	}
	if st.Round != 1 {
		t.Fatalf("round = %d, want 1", st.Round)
	}
	if got, want := len(st.Sequence), 2; got != want {
		t.Fatalf("sequence length = %d, want %d", got, want)
	}
	wantSequence := []int{2, 2}
	for i, want := range wantSequence {
		if st.Sequence[i] != want {
			t.Fatalf("sequence[%d] = %d, want %d", i, st.Sequence[i], want)
		}
	}
	if st.LitPad != wantSequence[0] {
		t.Fatalf("lit pad = %d, want %d", st.LitPad, wantSequence[0])
	}
}

func TestShowPhaseTransitionsToInput(t *testing.T) {
	t.Parallel()

	st := newState(t)
	if err := st.Step(game.Input{Start: true}); err != nil {
		t.Fatalf("Step start: %v", err)
	}

	// Two sequence items with show/gap timing should eventually finish the show phase.
	stepMany(t, st, 6)

	if st.Phase != game.PhaseInput {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseInput)
	}
	if st.InputIndex != 0 {
		t.Fatalf("input index = %d, want 0", st.InputIndex)
	}
	if st.LitPad != -1 {
		t.Fatalf("lit pad = %d, want -1 in input phase", st.LitPad)
	}
}

func TestCorrectSequenceAdvancesRoundAndShortensTiming(t *testing.T) {
	t.Parallel()

	st := newState(t)
	initialShowTicks := st.ShowTicks
	initialGapTicks := st.GapTicks

	if err := st.Step(game.Input{Start: true}); err != nil {
		t.Fatalf("Step start: %v", err)
	}
	stepMany(t, st, 6)

	for _, pad := range st.Sequence {
		if err := st.Step(game.Input{Pad: pad, Press: true}); err != nil {
			t.Fatalf("Step correct pad %d: %v", pad, err)
		}
	}

	if st.Phase != game.PhaseShowing {
		t.Fatalf("phase = %v, want %v after clearing round", st.Phase, game.PhaseShowing)
	}
	if st.Round != 2 {
		t.Fatalf("round = %d, want 2", st.Round)
	}
	if got, want := len(st.Sequence), 3; got != want {
		t.Fatalf("sequence length = %d, want %d", got, want)
	}
	if st.ShowTicks >= initialShowTicks {
		t.Fatalf("show ticks = %d, want < %d after ramp", st.ShowTicks, initialShowTicks)
	}
	if st.GapTicks >= initialGapTicks {
		t.Fatalf("gap ticks = %d, want < %d after ramp", st.GapTicks, initialGapTicks)
	}
}

func TestWrongInputEndsGame(t *testing.T) {
	t.Parallel()

	st := newState(t)
	if err := st.Step(game.Input{Start: true}); err != nil {
		t.Fatalf("Step start: %v", err)
	}
	stepMany(t, st, 6)

	if err := st.Step(game.Input{Pad: 0, Press: true}); err != nil {
		t.Fatalf("Step wrong pad: %v", err)
	}
	if st.Phase != game.PhaseGameOver {
		t.Fatalf("phase = %v, want %v after wrong input", st.Phase, game.PhaseGameOver)
	}
}

func TestWatchdogTimeoutEndsGame(t *testing.T) {
	t.Parallel()

	st := newState(t)
	if err := st.Step(game.Input{Start: true}); err != nil {
		t.Fatalf("Step start: %v", err)
	}
	stepMany(t, st, 6)
	stepMany(t, st, st.InputTimeoutTicks)

	if st.Phase != game.PhaseGameOver {
		t.Fatalf("phase = %v, want %v after watchdog timeout", st.Phase, game.PhaseGameOver)
	}
}
