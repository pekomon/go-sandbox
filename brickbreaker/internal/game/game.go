package game

import (
	"errors"
	"math"
)

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

func (s *State) Step(in Input) error {
	if s == nil {
		return errors.New("nil state")
	}

	switch s.Phase {
	case PhaseServe:
		s.stepServe(in)
		if s.Phase == PhaseServe {
			return nil
		}
	case PhaseWon, PhaseGameOver:
		return nil
	}

	prevX := s.Ball.X
	prevY := s.Ball.Y

	s.Ball.X += s.Ball.VX
	s.Ball.Y += s.Ball.VY

	s.bounceWalls()

	if s.Ball.Y+s.Ball.Radius > s.Height {
		s.consumeLife()
		return nil
	}

	s.bouncePaddle(prevX, prevY)
	s.hitBrick(prevX, prevY)

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

func (s *State) stepServe(in Input) {
	s.movePaddle(in)
	s.attachBallToPaddle()
	if !in.Launch {
		return
	}

	s.Phase = PhaseRunning
	s.Ball.VX = s.Ball.Speed * 0.6
	s.Ball.VY = -s.Ball.Speed * 0.8
}

func (s *State) movePaddle(in Input) {
	if in.MoveLeft == in.MoveRight {
		return
	}

	if in.MoveLeft {
		s.Paddle.X -= s.Paddle.Speed
	}
	if in.MoveRight {
		s.Paddle.X += s.Paddle.Speed
	}

	if s.Paddle.X < 0 {
		s.Paddle.X = 0
	}
	maxX := s.Width - s.Paddle.Width
	if s.Paddle.X > maxX {
		s.Paddle.X = maxX
	}
}

func (s *State) bounceWalls() {
	if s.Ball.X-s.Ball.Radius < 0 {
		s.Ball.X = s.Ball.Radius
		s.Ball.VX = math.Abs(s.Ball.VX)
	}
	if s.Ball.X+s.Ball.Radius > s.Width {
		s.Ball.X = s.Width - s.Ball.Radius
		s.Ball.VX = -math.Abs(s.Ball.VX)
	}
	if s.Ball.Y-s.Ball.Radius < 0 {
		s.Ball.Y = s.Ball.Radius
		s.Ball.VY = math.Abs(s.Ball.VY)
	}
}

func (s *State) consumeLife() {
	s.Lives--
	if s.Lives <= 0 {
		s.Lives = 0
		s.Phase = PhaseGameOver
		s.Ball.VX = 0
		s.Ball.VY = 0
		return
	}

	s.Phase = PhaseServe
	s.attachBallToPaddle()
}

func (s *State) bouncePaddle(_, prevY float64) {
	if s.Ball.VY <= 0 {
		return
	}

	paddleTop := s.Paddle.Y
	prevBottom := prevY + s.Ball.Radius
	ballBottom := s.Ball.Y + s.Ball.Radius
	if prevBottom > paddleTop || ballBottom < paddleTop {
		return
	}
	if s.Ball.X+s.Ball.Radius < s.Paddle.X || s.Ball.X-s.Ball.Radius > s.Paddle.X+s.Paddle.Width {
		return
	}

	s.Ball.Y = paddleTop - s.Ball.Radius
	s.Ball.VY = -math.Abs(s.Ball.Speed)

	relative := (s.Ball.X - (s.Paddle.X + s.Paddle.Width/2)) / (s.Paddle.Width / 2)
	if relative < -1 {
		relative = -1
	}
	if relative > 1 {
		relative = 1
	}
	s.Ball.VX = relative * s.Ball.Speed
	if math.Abs(s.Ball.VX) < s.Ball.Speed*0.2 {
		if relative < 0 {
			s.Ball.VX = -s.Ball.Speed * 0.2
		} else {
			s.Ball.VX = s.Ball.Speed * 0.2
		}
	}
}

func (s *State) hitBrick(prevX, prevY float64) {
	for i := range s.Bricks {
		if !s.Bricks[i].Alive {
			continue
		}
		if !ballIntersectsBrick(s.Ball, s.Bricks[i]) {
			continue
		}

		brick := &s.Bricks[i]
		brick.Alive = false
		s.Score++
		s.reflectFromBrick(*brick, prevX, prevY)

		if s.RemainingBricks() == 0 {
			s.Phase = PhaseWon
		}
		return
	}
}

func (s *State) reflectFromBrick(brick Brick, prevX, prevY float64) {
	prevLeft := prevX - s.Ball.Radius
	prevRight := prevX + s.Ball.Radius
	prevTop := prevY - s.Ball.Radius
	prevBottom := prevY + s.Ball.Radius

	brickLeft := brick.X
	brickRight := brick.X + brick.Width
	brickTop := brick.Y
	brickBottom := brick.Y + brick.Height

	switch {
	case prevBottom <= brickTop:
		s.Ball.Y = brickTop - s.Ball.Radius
		s.Ball.VY = -math.Abs(s.Ball.VY)
	case prevTop >= brickBottom:
		s.Ball.Y = brickBottom + s.Ball.Radius
		s.Ball.VY = math.Abs(s.Ball.VY)
	case prevRight <= brickLeft:
		s.Ball.X = brickLeft - s.Ball.Radius
		s.Ball.VX = -math.Abs(s.Ball.VX)
	case prevLeft >= brickRight:
		s.Ball.X = brickRight + s.Ball.Radius
		s.Ball.VX = math.Abs(s.Ball.VX)
	default:
		overlapLeft := math.Abs((s.Ball.X + s.Ball.Radius) - brickLeft)
		overlapRight := math.Abs(brickRight - (s.Ball.X - s.Ball.Radius))
		overlapTop := math.Abs((s.Ball.Y + s.Ball.Radius) - brickTop)
		overlapBottom := math.Abs(brickBottom - (s.Ball.Y - s.Ball.Radius))

		minOverlap := overlapLeft
		axis := "left"
		if overlapRight < minOverlap {
			minOverlap = overlapRight
			axis = "right"
		}
		if overlapTop < minOverlap {
			minOverlap = overlapTop
			axis = "top"
		}
		if overlapBottom < minOverlap {
			axis = "bottom"
		}

		switch axis {
		case "left":
			s.Ball.X = brickLeft - s.Ball.Radius
			s.Ball.VX = -math.Abs(s.Ball.VX)
		case "right":
			s.Ball.X = brickRight + s.Ball.Radius
			s.Ball.VX = math.Abs(s.Ball.VX)
		case "top":
			s.Ball.Y = brickTop - s.Ball.Radius
			s.Ball.VY = -math.Abs(s.Ball.VY)
		default:
			s.Ball.Y = brickBottom + s.Ball.Radius
			s.Ball.VY = math.Abs(s.Ball.VY)
		}
	}
}

func ballIntersectsBrick(ball Ball, brick Brick) bool {
	if ball.X+ball.Radius < brick.X || ball.X-ball.Radius > brick.X+brick.Width {
		return false
	}
	if ball.Y+ball.Radius < brick.Y || ball.Y-ball.Radius > brick.Y+brick.Height {
		return false
	}
	return true
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
