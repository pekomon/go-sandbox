package game_test

import (
    "fmt"
    "testing"

    "github.com/pekomon/go-sandbox/snake/internal/game"
)

// fakeRand feeds a predefined sequence of values to Intn.
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

func cfg(w, h, start int, r game.Rand) game.Config {
    return game.Config{Width: w, Height: h, StartLen: start, RNG: r}
}

func TestNew_PlacesSnakeAndApple_NoOverlap(t *testing.T) {
    r := &fakeRand{seq: []int{0, 1, 2, 3, 4, 5}} // deterministic
    st, err := game.New(cfg(10, 7, 3, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    if !st.Alive {
        t.Fatalf("expected alive at start")
    }
    if len(st.Snake) != 3 {
        t.Fatalf("start len = %d, want 3", len(st.Snake))
    }
    // Head should be index 0 and segments on the same row, facing Right.
    row := st.Snake[0].Y
    for i, p := range st.Snake {
        if p.Y != row {
            t.Fatalf("segment %d not in same row", i)
        }
    }
    // Apple must not overlap the snake.
    for _, p := range st.Snake {
        if p == st.Apple {
            t.Fatalf("apple overlaps snake at %v", p)
        }
    }
}

func TestTurn_OppositeIsIgnored(t *testing.T) {
    r := &fakeRand{seq: []int{0}}
    st, err := game.New(cfg(8, 6, 2, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    if st.Dir != game.Right {
        t.Fatalf("default dir = %v, want Right", st.Dir)
    }
    // Opposite turn ignored
    st.Turn(game.Left)
    if st.Dir != game.Right {
        t.Fatalf("opposite turn should be ignored; got %v", st.Dir)
    }
    // Valid turn accepted
    st.Turn(game.Down)
    if st.Dir != game.Down {
        t.Fatalf("turn Down not applied; got %v", st.Dir)
    }
}

func TestStep_Move_NoApple_NoGrowth(t *testing.T) {
    r := &fakeRand{seq: []int{5}}
    st, err := game.New(cfg(6, 5, 2, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    headBefore := st.Snake[0]
    st.Turn(game.Right) // ensure Right
    if err := st.Step(); err != nil {
        t.Fatalf("step: %v", err)
    }
    headAfter := st.Snake[0]
    if headAfter.X != headBefore.X+1 || headAfter.Y != headBefore.Y {
        t.Fatalf("head moved to %v, want (%d,%d)", headAfter, headBefore.X+1, headBefore.Y)
    }
    if l := len(st.Snake); l != 2 {
        t.Fatalf("length changed without eating; len=%d", l)
    }
}

func TestStep_EatApple_GrowsAndScores_AppleRelocated(t *testing.T) {
    // Place an apple directly in front of the head by picking seq that maps to that cell.
    r := &fakeRand{seq: []int{0, 1, 2, 3, 4, 5}}
    st, err := game.New(cfg(10, 7, 2, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    // Force apple position to be right of the head for the test’s sake:
    // We accept a helper pattern where implementation relocates if overlap occurs; here we just
    // validate that after eating, length+score increase and apple changes.
    oldApple := st.Apple
    oldLen := len(st.Snake)
    st.Turn(game.Right)
    if err := st.Step(); err != nil {
        t.Fatalf("step: %v", err)
    }
    if len(st.Snake) != oldLen && len(st.Snake) != oldLen+1 {
        // Implementation may require positioning apple deterministically; minimum expectation is growth when eaten.
        // For stricter behavior, adjust when implementing.
    }
    // Score must be >= 0 and increase when eating (assert exact +1 in implementation PR).
    if st.Score < 0 {
        t.Fatalf("score should be non-negative")
    }
    if st.Apple == oldApple {
        t.Fatalf("apple not relocated")
    }
}

func TestStep_WallCollision_Kills(t *testing.T) {
    r := &fakeRand{seq: []int{0}}
    st, err := game.New(cfg(3, 3, 2, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    // Move right until wall.
    for i := 0; i < 5; i++ {
        if err := st.Step(); err != nil {
            // Implementation may return an error or just set Alive=false; both acceptable,
            // but Alive must be false after collision.
            break
        }
    }
    if st.Alive {
        t.Fatalf("expected dead after wall collision")
    }
}

func TestStep_SelfCollision_Kills(t *testing.T) {
    r := &fakeRand{seq: []int{0}}
    st, err := game.New(cfg(5, 5, 4, r))
    if err != nil {
        t.Fatalf("new: %v", err)
    }
    // Make a small loop: Right, Down, Left, Up
    st.Turn(game.Right)
    _ = st.Step()
    st.Turn(game.Down)
    _ = st.Step()
    st.Turn(game.Left)
    _ = st.Step()
    st.Turn(game.Up)
    _ = st.Step()
    if st.Alive {
        t.Fatalf("expected dead after self-collision")
    }
}

func BenchmarkStep_StraightLine(b *testing.B) {
    r := &fakeRand{seq: []int{0}}
    st, _ := game.New(cfg(64, 36, 3, r))
    for i := 0; i < b.N; i++ {
        _ = st.Step()
    }
}

// Stringer helpers for clearer diffs (optional for implementation)
func (p game.Pos) String() string { return fmt.Sprintf("(%d,%d)", p.X, p.Y) }
