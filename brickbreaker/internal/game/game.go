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
	Levels [][]Brick
}

type State struct {
	Width       float64
	Height      float64
	Paddle      Paddle
	Ball        Ball
	Bricks      []Brick
	Lives       int
	Score       int
	Phase       Phase
	Level       int
	TotalLevels int

	initialPaddle  Paddle
	levelLayouts   [][]Brick
	baseBallSpeed  float64
	currentLevelID int
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
		Levels: defaultLevels(),
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
	layouts := cloneLevels(cfg.Levels)
	if len(layouts) == 0 {
		if len(cfg.Bricks) == 0 {
			cfg.Bricks = defaultBricks()
		}
		layouts = [][]Brick{cloneBricks(cfg.Bricks)}
	}
	st := &State{
		Width:         cfg.Width,
		Height:        cfg.Height,
		Paddle:        cfg.Paddle,
		Ball:          cfg.Ball,
		Lives:         cfg.Lives,
		Phase:         PhaseServe,
		Level:         1,
		TotalLevels:   len(layouts),
		initialPaddle: cfg.Paddle,
		levelLayouts:  layouts,
		baseBallSpeed: cfg.Ball.Speed,
	}
	st.loadLevel(0)
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
	case PhaseRunning:
		s.movePaddle(in)
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
	s.launchBall()
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
	speed := s.currentBallSpeed()
	relative := (s.Ball.X - (s.Paddle.X + s.Paddle.Width/2)) / (s.Paddle.Width / 2)
	if relative < -1 {
		relative = -1
	}
	if relative > 1 {
		relative = 1
	}
	if math.Abs(relative) < 0.2 {
		if relative < 0 {
			relative = -0.2
		} else {
			relative = 0.2
		}
	}
	angle := relative * (math.Pi / 3)
	s.Ball.VX = speed * math.Sin(angle)
	s.Ball.VY = -speed * math.Cos(angle)
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
			s.advanceLevel()
		}
		return
	}
}

func (s *State) loadLevel(levelIndex int) {
	s.currentLevelID = levelIndex
	s.Level = levelIndex + 1
	s.Paddle = s.initialPaddle
	s.Ball.Speed = s.levelBallSpeed(levelIndex)
	s.Bricks = cloneBricks(s.levelLayouts[levelIndex])
	s.Phase = PhaseServe
	s.attachBallToPaddle()
}

func (s *State) levelBallSpeed(levelIndex int) float64 {
	return s.baseBallSpeed + float64(levelIndex)*0.9
}

func (s *State) launchBall() {
	speed := s.currentBallSpeed()
	s.Ball.VX = speed * 0.6
	s.Ball.VY = -math.Sqrt(speed*speed - s.Ball.VX*s.Ball.VX)
}

func (s *State) currentBallSpeed() float64 {
	if s.Ball.Speed > 0 {
		return s.Ball.Speed
	}
	speed := math.Hypot(s.Ball.VX, s.Ball.VY)
	if speed == 0 {
		return 1
	}
	return speed
}

func (s *State) advanceLevel() {
	next := s.currentLevelID + 1
	if next >= s.TotalLevels {
		s.Phase = PhaseWon
		return
	}
	s.loadLevel(next)
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

func cloneLevels(levels [][]Brick) [][]Brick {
	if len(levels) == 0 {
		return nil
	}
	cloned := make([][]Brick, len(levels))
	for i := range levels {
		cloned[i] = cloneBricks(levels[i])
	}
	return cloned
}

func cloneBricks(bricks []Brick) []Brick {
	cloned := make([]Brick, len(bricks))
	copy(cloned, bricks)
	return cloned
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

func defaultLevels() [][]Brick {
	return [][]Brick{
		defaultBricks(),
		buildLevelLayout([]int{6, 7, 6, 7, 6}, 89, 52, 62, 24, 56, 16, true),
		buildLevelLayout([]int{4, 6, 8, 6, 4, 2}, 56, 48, 66, 22, 60, 16, false),
	}
}

func buildLevelLayout(rowWidths []int, startX, startY, stepX, stepY, brickWidth, brickHeight float64, stagger bool) []Brick {
	total := 0
	for _, count := range rowWidths {
		total += count
	}
	bricks := make([]Brick, 0, total)
	for row, count := range rowWidths {
		rowStartX := startX
		if stagger && row%2 == 1 {
			rowStartX += stepX / 2
		}
		for col := 0; col < count; col++ {
			bricks = append(bricks, Brick{
				X:      rowStartX + float64(col)*stepX,
				Y:      startY + float64(row)*stepY,
				Width:  brickWidth,
				Height: brickHeight,
				Alive:  true,
			})
		}
	}
	return bricks
}
