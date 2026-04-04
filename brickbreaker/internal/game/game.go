package game

import "errors"

type Phase int

const (
	PhaseServe Phase = iota
	PhaseRunning
	PhaseWon
	PhaseGameOver
)

type Input struct {
	MoveLeft  bool
	MoveRight bool
	Launch    bool
}

type Paddle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Speed  float64
}

type Ball struct {
	X      float64
	Y      float64
	VX     float64
	VY     float64
	Radius float64
	Speed  float64
}

type Brick struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Alive  bool
}

type Config struct {
	Width  float64
	Height float64
	Paddle Paddle
	Ball   Ball
	Lives  int
	Bricks []Brick
}

type State struct {
	Width  float64
	Height float64
	Paddle Paddle
	Ball   Ball
	Bricks []Brick
	Lives  int
	Score  int
	Phase  Phase
}

func DefaultConfig() Config {
	return Config{
		Width:  640,
		Height: 480,
		Paddle: Paddle{
			X:      272,
			Y:      440,
			Width:  96,
			Height: 16,
			Speed:  12,
		},
		Ball: Ball{
			Radius: 8,
			Speed:  6,
		},
		Lives:  3,
		Bricks: defaultBricks(),
	}
}

func New(cfg Config) (*State, error) {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, errors.New("invalid dimensions")
	}
	if cfg.Paddle.Width <= 0 || cfg.Paddle.Height <= 0 {
		return nil, errors.New("invalid paddle")
	}
	if cfg.Ball.Radius <= 0 {
		return nil, errors.New("invalid ball")
	}
	if cfg.Lives <= 0 {
		return nil, errors.New("invalid lives")
	}
	bricks := make([]Brick, len(cfg.Bricks))
	copy(bricks, cfg.Bricks)
	st := &State{
		Width:  cfg.Width,
		Height: cfg.Height,
		Paddle: cfg.Paddle,
		Ball:   cfg.Ball,
		Bricks: bricks,
		Lives:  cfg.Lives,
		Phase:  PhaseServe,
	}
	st.attachBallToPaddle()
	return st, nil
}

func (s *State) Step(Input) error {
	return nil
}

func (s *State) RemainingBricks() int {
	count := 0
	for _, brick := range s.Bricks {
		if brick.Alive {
			count++
		}
	}
	return count
}

func (s *State) attachBallToPaddle() {
	s.Ball.X = s.Paddle.X + s.Paddle.Width/2
	s.Ball.Y = s.Paddle.Y - s.Ball.Radius - 1
	s.Ball.VX = 0
	s.Ball.VY = 0
}

func defaultBricks() []Brick {
	bricks := make([]Brick, 0, 40)
	for row := 0; row < 5; row++ {
		for col := 0; col < 8; col++ {
			bricks = append(bricks, Brick{
				X:      56 + float64(col)*66,
				Y:      52 + float64(row)*26,
				Width:  60,
				Height: 18,
				Alive:  true,
			})
		}
	}
	return bricks
}
