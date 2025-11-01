package game_test

import (
	"errors"
	"testing"

	"github.com/pekomon/go-sandbox/snake/internal/game"
)

type fakeRand struct {
	seq []int
	i   int
}

func (f *fakeRand) Intn(n int) int {
	if len(f.seq) == 0 {
		return 0
	}
	v := f.seq[f.i%len(f.seq)] % n
	f.i++
	return v
}

func newConfig(w, h, start int, r game.Rand) game.Config {
	return game.Config{Width: w, Height: h, StartLen: start, RNG: r}
}

func TestNewInitialState(t *testing.T) {
	t.Parallel()

	r := &fakeRand{seq: []int{3, 2, 4, 1}}
	st, err := game.New(newConfig(10, 7, 3, r))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	if !st.Alive {
		t.Fatalf("expected Alive=true at start")
	}
	if len(st.Snake) != 3 {
		t.Fatalf("start length = %d, want 3", len(st.Snake))
	}
	head := st.Snake[0]
	for i, seg := range st.Snake {
		if seg.Y != head.Y {
			t.Fatalf("segment %d not on same row as head (row=%d, got %d)", i, head.Y, seg.Y)
		}
	}
	for _, seg := range st.Snake {
		if seg == st.Apple {
			t.Fatalf("apple overlaps snake segment %v", seg)
		}
	}
}

func TestTurnIgnoresOpposite(t *testing.T) {
	t.Parallel()

	st, err := game.New(newConfig(8, 6, 2, &fakeRand{}))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	if st.Dir != game.Right {
		t.Fatalf("default direction = %v, want Right", st.Dir)
	}

	st.Turn(game.Left)
	if st.Dir != game.Right {
		t.Fatalf("opposite direction should be ignored; got %v", st.Dir)
	}

	st.Turn(game.Down)
	if st.Dir != game.Down {
		t.Fatalf("expected Down after Turn; got %v", st.Dir)
	}
}

func TestStepMovesSnake(t *testing.T) {
	t.Parallel()

	st, err := game.New(newConfig(6, 5, 2, &fakeRand{seq: []int{4, 3, 2, 1}}))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	headBefore := st.Snake[0]
	if err := st.Step(); err != nil {
		t.Fatalf("Step: %v", err)
	}
	headAfter := st.Snake[0]
	if headAfter.X != headBefore.X+1 || headAfter.Y != headBefore.Y {
		t.Fatalf("head moved to %v, want (%d,%d)", headAfter, headBefore.X+1, headBefore.Y)
	}
	if got, want := len(st.Snake), 2; got != want {
		t.Fatalf("length changed without apple; got %d want %d", got, want)
	}
}

func TestStepEatAppleIncreasesScoreAndLength(t *testing.T) {
	t.Parallel()

	st, err := game.New(newConfig(8, 6, 2, &fakeRand{seq: []int{5, 2, 6, 4}}))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	head := st.Snake[0]
	st.Apple = game.Pos{X: head.X + 1, Y: head.Y}

	oldLen := len(st.Snake)
	oldScore := st.Score
	oldApple := st.Apple

	if err := st.Step(); err != nil {
		t.Fatalf("Step: %v", err)
	}
	if got, want := len(st.Snake), oldLen+1; got != want {
		t.Fatalf("snake length = %d, want %d", got, want)
	}
	if got, want := st.Score, oldScore+1; got != want {
		t.Fatalf("score = %d, want %d", got, want)
	}
	if st.Apple == oldApple {
		t.Fatalf("apple not relocated")
	}
}

func TestStepWallCollisionStopsGame(t *testing.T) {
	t.Parallel()

	st, err := game.New(newConfig(3, 3, 2, &fakeRand{}))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}

	var collision error
	for i := 0; i < 5; i++ {
		collision = st.Step()
		if errors.Is(collision, game.ErrGameOver) {
			break
		}
	}
	if !errors.Is(collision, game.ErrGameOver) {
		t.Fatalf("expected ErrGameOver, got %v", collision)
	}
	if st.Alive {
		t.Fatalf("Alive should be false after wall collision")
	}
}

func TestStepSelfCollisionStopsGame(t *testing.T) {
	t.Parallel()

	st, err := game.New(newConfig(6, 6, 3, &fakeRand{}))
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	st.Snake = []game.Pos{
		{X: 3, Y: 3},
		{X: 2, Y: 3},
		{X: 2, Y: 2},
		{X: 3, Y: 2},
	}
	st.Apple = game.Pos{X: 0, Y: 0}
	st.Dir = game.Left

	if err := st.Step(); !errors.Is(err, game.ErrGameOver) {
		t.Fatalf("expected ErrGameOver on self-collision, got %v", err)
	}
	if st.Alive {
		t.Fatalf("Alive should be false after self-collision")
	}
}
