package game_test

import (
	"math"
	"testing"

	"github.com/pekomon/go-sandbox/brickbreaker/internal/game"
)

func newState(t *testing.T) *game.State {
	t.Helper()

	st, err := game.New(game.DefaultConfig())
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	return st
}

func newStateWithConfig(t *testing.T, cfg game.Config) *game.State {
	t.Helper()

	st, err := game.New(cfg)
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	return st
}

func TestNewStartsInServeState(t *testing.T) {
	t.Parallel()

	st := newState(t)

	if st.Phase != game.PhaseServe {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseServe)
	}
	if st.Lives != 3 {
		t.Fatalf("lives = %d, want 3", st.Lives)
	}
	if st.Level != 1 {
		t.Fatalf("level = %d, want 1", st.Level)
	}
	if st.TotalLevels < 3 {
		t.Fatalf("total levels = %d, want at least 3", st.TotalLevels)
	}
	if st.RemainingBricks() == 0 {
		t.Fatalf("expected a populated brick field at start")
	}
	wantBallX := st.Paddle.X + st.Paddle.Width/2
	wantBallY := st.Paddle.Y - st.Ball.Radius - 1
	if st.Ball.X != wantBallX || st.Ball.Y != wantBallY {
		t.Fatalf("ball start = (%.1f, %.1f), want (%.1f, %.1f)", st.Ball.X, st.Ball.Y, wantBallX, wantBallY)
	}
}

func TestClearingIntermediateLevelAdvancesToNextServe(t *testing.T) {
	t.Parallel()

	cfg := game.DefaultConfig()
	cfg.Levels = [][]game.Brick{
		{{
			X:      160,
			Y:      80,
			Width:  60,
			Height: 18,
			Alive:  true,
		}},
		{
			{X: 100, Y: 60, Width: 60, Height: 18, Alive: true},
			{X: 200, Y: 60, Width: 60, Height: 18, Alive: true},
		},
	}
	cfg.Bricks = nil

	st, err := game.New(cfg)
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}
	initialSpeed := st.Ball.Speed

	st.Phase = game.PhaseRunning
	st.Ball.X = st.Bricks[0].X + st.Bricks[0].Width/2
	st.Ball.Y = st.Bricks[0].Y + st.Bricks[0].Height + st.Ball.Radius + 1
	st.Ball.VX = 0
	st.Ball.VY = -4

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step level clear: %v", err)
	}
	if st.Phase != game.PhaseServe {
		t.Fatalf("phase = %v, want %v after progression", st.Phase, game.PhaseServe)
	}
	if st.Level != 2 {
		t.Fatalf("level = %d, want 2", st.Level)
	}
	if st.TotalLevels != 2 {
		t.Fatalf("total levels = %d, want 2", st.TotalLevels)
	}
	if st.RemainingBricks() != 2 {
		t.Fatalf("remaining bricks = %d, want 2 from next layout", st.RemainingBricks())
	}
	if st.Ball.Speed <= initialSpeed {
		t.Fatalf("ball speed = %.2f, want > %.2f on next level", st.Ball.Speed, initialSpeed)
	}
}

func TestLaunchUsesCurrentLevelSpeed(t *testing.T) {
	t.Parallel()

	cfg := game.DefaultConfig()
	cfg.Levels = [][]game.Brick{
		{{X: 160, Y: 80, Width: 60, Height: 18, Alive: true}},
		{{X: 160, Y: 80, Width: 60, Height: 18, Alive: true}},
	}
	cfg.Bricks = nil

	st, err := game.New(cfg)
	if err != nil {
		t.Fatalf("game.New: %v", err)
	}

	st.Phase = game.PhaseRunning
	st.Ball.X = st.Bricks[0].X + st.Bricks[0].Width/2
	st.Ball.Y = st.Bricks[0].Y + st.Bricks[0].Height + st.Ball.Radius + 1
	st.Ball.VX = 0
	st.Ball.VY = -4

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step level clear: %v", err)
	}
	if st.Level != 2 {
		t.Fatalf("level = %d, want 2", st.Level)
	}

	expectedSpeed := st.Ball.Speed
	if err := st.Step(game.Input{Launch: true}); err != nil {
		t.Fatalf("Step launch after progression: %v", err)
	}
	gotSpeed := math.Hypot(st.Ball.VX, st.Ball.VY)
	if math.Abs(gotSpeed-expectedSpeed) > 0.001 {
		t.Fatalf("launch speed = %.3f, want %.3f", gotSpeed, expectedSpeed)
	}
}

func TestStepServeMovementMovesPaddleAndBall(t *testing.T) {
	t.Parallel()

	st := newState(t)
	startPaddleX := st.Paddle.X
	startBallX := st.Ball.X

	if err := st.Step(game.Input{MoveRight: true}); err != nil {
		t.Fatalf("Step: %v", err)
	}
	if st.Paddle.X <= startPaddleX {
		t.Fatalf("paddle X = %.1f, want > %.1f", st.Paddle.X, startPaddleX)
	}
	if st.Ball.X <= startBallX {
		t.Fatalf("ball X = %.1f, want > %.1f while serving", st.Ball.X, startBallX)
	}

	st.Paddle.X = st.Width - st.Paddle.Width
	st.Ball.X = st.Paddle.X + st.Paddle.Width/2
	if err := st.Step(game.Input{MoveRight: true}); err != nil {
		t.Fatalf("Step at boundary: %v", err)
	}
	if st.Paddle.X != st.Width-st.Paddle.Width {
		t.Fatalf("paddle escaped right wall: got %.1f", st.Paddle.X)
	}
}

func TestStepLaunchStartsBallMotion(t *testing.T) {
	t.Parallel()

	st := newState(t)

	if err := st.Step(game.Input{Launch: true}); err != nil {
		t.Fatalf("Step launch: %v", err)
	}
	if st.Phase != game.PhaseRunning {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseRunning)
	}
	if st.Ball.VY >= 0 {
		t.Fatalf("ball VY = %.1f, want upward launch", st.Ball.VY)
	}
	if st.Ball.VX == 0 {
		t.Fatalf("ball VX = %.1f, want non-zero horizontal launch", st.Ball.VX)
	}
}

func TestStepWallBounceReflectsVelocity(t *testing.T) {
	t.Parallel()

	st := newState(t)
	st.Phase = game.PhaseRunning
	st.Ball.X = st.Ball.Radius + 1
	st.Ball.Y = st.Ball.Radius + 1
	st.Ball.VX = -4
	st.Ball.VY = -3

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step wall bounce: %v", err)
	}
	if st.Ball.VX <= 0 {
		t.Fatalf("ball VX = %.1f, want reflected horizontal velocity", st.Ball.VX)
	}
	if st.Ball.VY <= 0 {
		t.Fatalf("ball VY = %.1f, want reflected vertical velocity from top wall", st.Ball.VY)
	}
}

func TestStepPaddleBounceReflectsBallUpward(t *testing.T) {
	t.Parallel()

	st := newState(t)
	st.Phase = game.PhaseRunning
	st.Ball.X = st.Paddle.X + st.Paddle.Width/2
	st.Ball.Y = st.Paddle.Y - st.Ball.Radius - 1
	st.Ball.VX = 0
	st.Ball.VY = 5

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step paddle bounce: %v", err)
	}
	if st.Ball.VY >= 0 {
		t.Fatalf("ball VY = %.1f, want upward bounce from paddle", st.Ball.VY)
	}
}

func TestStepRunningMovementMovesPaddle(t *testing.T) {
	t.Parallel()

	st := newState(t)
	st.Phase = game.PhaseRunning
	startPaddleX := st.Paddle.X

	if err := st.Step(game.Input{MoveRight: true}); err != nil {
		t.Fatalf("Step running move: %v", err)
	}
	if st.Paddle.X <= startPaddleX {
		t.Fatalf("paddle X = %.1f, want > %.1f while running", st.Paddle.X, startPaddleX)
	}

	st.Paddle.X = st.Width - st.Paddle.Width
	if err := st.Step(game.Input{MoveRight: true}); err != nil {
		t.Fatalf("Step running boundary: %v", err)
	}
	if st.Paddle.X != st.Width-st.Paddle.Width {
		t.Fatalf("paddle escaped right wall while running: got %.1f", st.Paddle.X)
	}
}

func TestStepBrickCollisionClearsBrickAndScores(t *testing.T) {
	t.Parallel()

	cfg := game.DefaultConfig()
	cfg.Levels = nil
	cfg.Bricks = []game.Brick{{
		X:      160,
		Y:      80,
		Width:  60,
		Height: 18,
		Alive:  true,
	}}
	st := newStateWithConfig(t, cfg)
	st.Phase = game.PhaseRunning
	st.Ball.X = st.Bricks[0].X + st.Bricks[0].Width/2
	st.Ball.Y = st.Bricks[0].Y + st.Bricks[0].Height + st.Ball.Radius + 1
	st.Ball.VX = 0
	st.Ball.VY = -4

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step brick collision: %v", err)
	}
	if st.Bricks[0].Alive {
		t.Fatalf("brick should be cleared after collision")
	}
	if st.Score != 1 {
		t.Fatalf("score = %d, want 1", st.Score)
	}
	if st.Ball.VY <= 0 {
		t.Fatalf("ball VY = %.1f, want downward reflection after hitting brick from below", st.Ball.VY)
	}
}

func TestStepBottomOutConsumesLifeAndResetsServe(t *testing.T) {
	t.Parallel()

	st := newState(t)
	st.Phase = game.PhaseRunning
	st.Lives = 3
	st.Ball.Y = st.Height - st.Ball.Radius - 1
	st.Ball.VY = 8

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step bottom out: %v", err)
	}
	if st.Lives != 2 {
		t.Fatalf("lives = %d, want 2", st.Lives)
	}
	if st.Phase != game.PhaseServe {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseServe)
	}
	wantBallX := st.Paddle.X + st.Paddle.Width/2
	wantBallY := st.Paddle.Y - st.Ball.Radius - 1
	if st.Ball.X != wantBallX || st.Ball.Y != wantBallY {
		t.Fatalf("ball reset = (%.1f, %.1f), want (%.1f, %.1f)", st.Ball.X, st.Ball.Y, wantBallX, wantBallY)
	}
}

func TestStepClearingLastBrickWinsRound(t *testing.T) {
	t.Parallel()

	cfg := game.DefaultConfig()
	cfg.Levels = nil
	cfg.Bricks = []game.Brick{{
		X:      200,
		Y:      80,
		Width:  60,
		Height: 18,
		Alive:  true,
	}}
	st := newStateWithConfig(t, cfg)
	st.Phase = game.PhaseRunning
	st.Ball.X = st.Bricks[0].X + st.Bricks[0].Width/2
	st.Ball.Y = st.Bricks[0].Y + st.Bricks[0].Height + st.Ball.Radius + 1
	st.Ball.VX = 0
	st.Ball.VY = -4

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step last brick: %v", err)
	}
	if st.Phase != game.PhaseWon {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseWon)
	}
	if st.RemainingBricks() != 0 {
		t.Fatalf("remaining bricks = %d, want 0", st.RemainingBricks())
	}
}

func TestStepBottomOutWithNoLivesSetsGameOver(t *testing.T) {
	t.Parallel()

	st := newState(t)
	st.Phase = game.PhaseRunning
	st.Lives = 1
	st.Ball.Y = st.Height - st.Ball.Radius - 1
	st.Ball.VY = 8

	if err := st.Step(game.Input{}); err != nil {
		t.Fatalf("Step game over: %v", err)
	}
	if st.Lives != 0 {
		t.Fatalf("lives = %d, want 0", st.Lives)
	}
	if st.Phase != game.PhaseGameOver {
		t.Fatalf("phase = %v, want %v", st.Phase, game.PhaseGameOver)
	}
}
