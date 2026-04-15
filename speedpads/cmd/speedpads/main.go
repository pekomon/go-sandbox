package main

import (
	"errors"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/pekomon/go-sandbox/speedpads/internal/game"
)

var backgroundColor = color.RGBA{0x10, 0x12, 0x18, 0xff}

type app struct {
	state *game.State
}

func newApp() (*app, error) {
	st, err := game.New(game.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &app{state: st}, nil
}

func (a *app) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Speedpads scaffold\n\nGameplay loop arrives in the follow-up feature PR.\n\nPlanned controls:\n- Space: start\n- D/F/J/K: input pads\n- R: restart after game over\n- Esc: quit")
}

func (a *app) Layout(int, int) (int, int) {
	return int(a.state.Width), int(a.state.Height)
}

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatalf("init speedpads: %v", err)
	}

	ebiten.SetWindowTitle("Speedpads")
	ebiten.SetWindowSize(int(app.state.Width), int(app.state.Height))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(app); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}
