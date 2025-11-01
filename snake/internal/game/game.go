package game

import "errors"

var ErrNotImplemented = errors.New("not implemented")

// Dir is the movement direction.
type Dir int

const (
    Up Dir = iota
    Right
    Down
    Left
)

// Opposite returns true if d is the opposite of other.
func (d Dir) Opposite(other Dir) bool {
    return (d == Up && other == Down) ||
        (d == Down && other == Up) ||
        (d == Left && other == Right) ||
        (d == Right && other == Left)
}

type Pos struct{ X, Y int }

// Rand abstracts randomness for deterministic tests.
type Rand interface {
    Intn(n int) int
}

// Config defines board size and initial snake length.
type Config struct {
    Width, Height int // required: >0
    StartLen      int // >= 1
    RNG           Rand // required for deterministic apple placement in tests
}

// State holds world state.
type State struct {
    Snake []Pos // head is index 0
    Dir   Dir
    Apple Pos
    Alive bool
    Score int

    w, h int
    rng  Rand
}

// New initializes a new State with snake centered horizontally, facing Right,
// and an apple placed in a free cell using RNG.
func New(cfg Config) (*State, error) {
    return nil, ErrNotImplemented
}

// Turn requests a direction change; reversing into the immediate opposite is ignored.
func (s *State) Turn(d Dir) {
    // STUB: not implemented
}

// Step advances the game by one tick: move head by Dir, shift body,
// detect wall/self collisions, handle eating (growth+score), and relocate apple.
func (s *State) Step() error {
    return ErrNotImplemented
}
