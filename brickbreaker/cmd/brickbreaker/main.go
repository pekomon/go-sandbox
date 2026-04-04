package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/pekomon/go-sandbox/brickbreaker/internal/game"
)

var (
	backgroundColor = color.RGBA{0x10, 0x18, 0x20, 0xff}
	brickColor      = color.RGBA{0xf1, 0xb8, 0x2d, 0xff}
	paddleColor     = color.RGBA{0x58, 0xb7, 0xe0, 0xff}
	ballColor       = color.RGBA{0xf7, 0xf3, 0xe8, 0xff}
)

type app struct {
	state *game.State
}

func newApp() (*app, error) {
	a := &app{}
	if err := a.reset(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *app) reset() error {
	st, err := game.New(game.DefaultConfig())
	if err != nil {
		return err
	}
	a.state = st
	return nil
}

func (a *app) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if (a.state.Phase == game.PhaseWon || a.state.Phase == game.PhaseGameOver) &&
		(inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		return a.reset()
	}

	in := game.Input{
		MoveLeft:  ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA),
		MoveRight: ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD),
		Launch:    inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	return a.state.Step(in)
}

func (a *app) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	for _, brick := range a.state.Bricks {
		if !brick.Alive {
			continue
		}
		ebitenutil.DrawRect(screen, brick.X, brick.Y, brick.Width, brick.Height, brickColor)
	}

	ebitenutil.DrawRect(screen, a.state.Paddle.X, a.state.Paddle.Y, a.state.Paddle.Width, a.state.Paddle.Height, paddleColor)
	ebitenutil.DrawRect(screen, a.state.Ball.X-a.state.Ball.Radius, a.state.Ball.Y-a.state.Ball.Radius, a.state.Ball.Radius*2, a.state.Ball.Radius*2, ballColor)

	msg := fmt.Sprintf("Lives: %d  Score: %d", a.state.Lives, a.state.Score)
	switch a.state.Phase {
	case game.PhaseServe:
		msg += "   Left/Right or A/D to move   Space to launch   Esc to quit"
	case game.PhaseRunning:
		speed := math.Hypot(a.state.Ball.VX, a.state.Ball.VY)
		msg += fmt.Sprintf("   Ball speed: %.1f   Esc to quit", speed)
	case game.PhaseWon:
		msg += "   YOU WIN   Press R, Space, or Enter to restart"
	case game.PhaseGameOver:
		msg += "   GAME OVER   Press R, Space, or Enter to restart"
	}
	ebitenutil.DebugPrintAt(screen, msg, 12, 12)
}

func (a *app) Layout(int, int) (int, int) {
	return int(a.state.Width), int(a.state.Height)
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatalf("init brickbreaker: %v", err)
	}

	ebiten.SetWindowTitle("BrickBreaker")
	ebiten.SetWindowSize(int(app.state.Width), int(app.state.Height))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(app); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}
