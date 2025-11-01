package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/pekomon/go-sandbox/snake/internal/game"
)

const (
	boardWidth      = 24
	boardHeight     = 18
	tileSize        = 24
	initialMoveWait = 6 // frames between steps (60 FPS default)
)

var (
	backgroundColor = color.RGBA{0x12, 0x16, 0x1c, 0xff}
	gridColor       = color.RGBA{0x1c, 0x23, 0x2b, 0xff}
	snakeColor      = color.RGBA{0x3c, 0xd0, 0x73, 0xff}
	headColor       = color.RGBA{0x70, 0xe0, 0x8a, 0xff}
	appleColor      = color.RGBA{0xe6, 0x53, 0x53, 0xff}
)

type snakeGame struct {
	state        *game.State
	frameCounter int
	moveDelay    int
	rng          *rand.Rand
}

func newSnakeGame() (*snakeGame, error) {
	g := &snakeGame{
		moveDelay: initialMoveWait,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if err := g.reset(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *snakeGame) reset() error {
	cfg := game.Config{
		Width:    boardWidth,
		Height:   boardHeight,
		StartLen: 3,
		RNG:      rand.New(rand.NewSource(g.rng.Int63())),
	}
	st, err := game.New(cfg)
	if err != nil {
		return err
	}
	g.state = st
	g.frameCounter = 0
	g.moveDelay = initialMoveWait
	return nil
}

func (g *snakeGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) || (!g.state.Alive && (inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter))) {
		return g.reset()
	}

	g.handleInput()

	if !g.state.Alive {
		return nil
	}

	g.frameCounter++
	if g.frameCounter < g.moveDelay {
		return nil
	}
	g.frameCounter = 0

	if err := g.state.Step(); err != nil {
		if errors.Is(err, game.ErrGameOver) {
			return nil
		}
		return err
	}

	// Speed up slightly every 5 apples, with a minimum delay.
	if delay := initialMoveWait - g.state.Score/5; delay >= 2 && delay < g.moveDelay {
		g.moveDelay = delay
	}
	return nil
}

func (g *snakeGame) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.state.Turn(game.Up)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.state.Turn(game.Down)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.state.Turn(game.Left)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.state.Turn(game.Right)
	}
}

func (g *snakeGame) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.drawGrid(screen)

	if g.state == nil {
		return
	}

	// Draw snake body.
	for i, seg := range g.state.Snake {
		c := snakeColor
		if i == 0 {
			c = headColor
		}
		ebitenutil.DrawRect(screen, float64(seg.X*tileSize), float64(seg.Y*tileSize), float64(tileSize), float64(tileSize), c)
	}

	// Draw apple.
	ebitenutil.DrawRect(screen, float64(g.state.Apple.X*tileSize), float64(g.state.Apple.Y*tileSize), float64(tileSize), float64(tileSize), appleColor)

	msg := fmt.Sprintf("Score: %d  Speed: %.1f/s   (Arrows/WASD to move)", g.state.Score, 60.0/float64(g.moveDelay))
	if !g.state.Alive {
		msg += "   GAME OVER — press R, Space, or Enter to restart"
	} else {
		msg += "   Press Esc to quit"
	}
	ebitenutil.DebugPrintAt(screen, msg, 8, 8)
}

func (g *snakeGame) drawGrid(screen *ebiten.Image) {
	w := boardWidth * tileSize
	h := boardHeight * tileSize
	for x := 0; x < w; x += tileSize {
		ebitenutil.DrawRect(screen, float64(x), 0, 1, float64(h), gridColor)
	}
	for y := 0; y < h; y += tileSize {
		ebitenutil.DrawRect(screen, 0, float64(y), float64(w), 1, gridColor)
	}
}

func (g *snakeGame) Layout(int, int) (int, int) {
	return boardWidth * tileSize, boardHeight * tileSize
}

func main() {
	ebiten.SetWindowTitle("Snake")
	ebiten.SetWindowSize(boardWidth*tileSize, boardHeight*tileSize)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	g, err := newSnakeGame()
	if err != nil {
		log.Fatalf("init snake: %v", err)
	}

	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatal(err)
	}
}
