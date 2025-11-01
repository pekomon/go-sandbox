package game

import (
	"errors"
	"math/rand"
)

var (
	ErrGameOver = errors.New("game over")
)

type Dir int

const (
	Up Dir = iota
	Right
	Down
	Left
)

func (d Dir) Opposite(other Dir) bool {
	return (d == Up && other == Down) ||
		(d == Down && other == Up) ||
		(d == Left && other == Right) ||
		(d == Right && other == Left)
}

type Pos struct{ X, Y int }

type Rand interface {
	Intn(n int) int
}

type Config struct {
	Width, Height int
	StartLen      int
	RNG           Rand
}

type State struct {
	Snake []Pos // head = index 0
	Dir   Dir
	Apple Pos
	Alive bool
	Score int

	w, h int
	rng  Rand
}

func New(cfg Config) (*State, error) {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, errors.New("invalid dimensions")
	}
	if cfg.StartLen < 1 {
		cfg.StartLen = 1
	}
	if cfg.StartLen > cfg.Width {
		cfg.StartLen = cfg.Width
	}
	if cfg.RNG == nil {
		cfg.RNG = rand.New(rand.NewSource(1))
	}
	s := &State{
		Dir:   Right,
		Alive: true,
		w:     cfg.Width,
		h:     cfg.Height,
		rng:   cfg.RNG,
	}
	cy := cfg.Height / 2
	cx := cfg.Width / 2
	startX := cx - (cfg.StartLen - 1)
	if startX < 0 {
		startX = 0
	}
	for i := 0; i < cfg.StartLen; i++ {
		x := startX + (cfg.StartLen - 1 - i)
		s.Snake = append(s.Snake, Pos{X: x, Y: cy})
	}
	s.placeApple()
	return s, nil
}

func (s *State) Turn(d Dir) {
	if d == s.Dir || d.Opposite(s.Dir) {
		return
	}
	s.Dir = d
}

func (s *State) Step() error {
	if !s.Alive {
		return ErrGameOver
	}
	if len(s.Snake) == 0 {
		return errors.New("empty snake")
	}
	head := s.Snake[0]
	next := head
	switch s.Dir {
	case Up:
		next.Y--
	case Down:
		next.Y++
	case Left:
		next.X--
	case Right:
		next.X++
	}

	if next.X < 0 || next.X >= s.w || next.Y < 0 || next.Y >= s.h {
		s.Alive = false
		return ErrGameOver
	}
	for _, p := range s.Snake {
		if p == next {
			s.Alive = false
			return ErrGameOver
		}
	}

	newSnake := append([]Pos{next}, s.Snake...)
	if next == s.Apple {
		s.Score++
		s.placeApple()
	} else {
		newSnake = newSnake[:len(newSnake)-1]
	}
	s.Snake = newSnake
	return nil
}

func (s *State) Size() (int, int) {
	return s.w, s.h
}

func (s *State) placeApple() {
	total := s.w * s.h
	for tries := 0; tries < total*2; tries++ {
		x := s.rng.Intn(s.w)
		y := s.rng.Intn(s.h)
		p := Pos{X: x, Y: y}
		if !s.occupied(p) {
			s.Apple = p
			return
		}
	}
	s.Apple = Pos{0, 0}
}

func (s *State) occupied(p Pos) bool {
	for _, q := range s.Snake {
		if q == p {
			return true
		}
	}
	return false
}
