package ui

import (
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/pekomon/go-sandbox/memesweeper/internal/board"
)

const (
	windowTitle = "MemeSweeper"
)

var (
	backgroundColor = color.RGBA{0x16, 0x1a, 0x20, 0xff}
	gridColor       = color.RGBA{0x22, 0x2b, 0x36, 0xff}
	hiddenColor     = color.RGBA{0x2c, 0x36, 0x45, 0xff}
	revealedColor   = color.RGBA{0x3e, 0x4a, 0x5c, 0xff}
)

// Config controls the game loop and board dimensions.
type Config struct {
	Preset    board.Preset
	Rows      int
	Cols      int
	MemeCount int
	TileSize  int
	Seed      int64
}

// Game implements the Ebiten game loop for MemeSweeper.
type Game struct {
	cfg        Config
	difficulty board.Preset
	assets     *Assets
	board      *board.Board
	rng        *rand.Rand
}

// Run loads assets, configures the window, and starts the Ebiten loop.
func Run(cfg Config) error {
	assets, err := LoadAssets()
	if err != nil {
		return err
	}
	game, err := NewGame(cfg, assets)
	if err != nil {
		return err
	}

	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetWindowSize(game.cfg.Cols*game.cfg.TileSize, game.cfg.Rows*game.cfg.TileSize)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game); err != nil && !errors.Is(err, ebiten.Termination) {
		return err
	}
	return nil
}

// NewGame builds a new Ebiten game instance.
func NewGame(cfg Config, assets *Assets) (*Game, error) {
	if cfg.TileSize <= 0 {
		return nil, fmt.Errorf("ui: invalid tile size")
	}
	if assets == nil || assets.Meme == nil || assets.Flag == nil {
		return nil, fmt.Errorf("ui: missing assets")
	}
	if cfg.Preset == "" && (cfg.Rows <= 0 || cfg.Cols <= 0) {
		return nil, fmt.Errorf("ui: invalid board dimensions")
	}

	g := &Game{
		cfg:        cfg,
		difficulty: cfg.Preset,
		assets:     assets,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if err := g.applyPreset(cfg.Preset); err != nil {
		return nil, err
	}
	if g.cfg.MemeCount < 0 || g.cfg.MemeCount > g.cfg.Rows*g.cfg.Cols {
		return nil, fmt.Errorf("ui: invalid meme count")
	}
	if err := g.reset(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.handleDifficultyShortcuts() {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) ||
		(g.board.Status != board.StatusActive &&
			(inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter))) {
		return g.reset()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return g.handleReveal()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.handleFlag()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	g.drawGrid(screen)

	for r := 0; r < g.cfg.Rows; r++ {
		for c := 0; c < g.cfg.Cols; c++ {
			g.drawCell(screen, r, c)
		}
	}

	flags := countFlags(g.board)
	status := statusLabel(g.board.Status)
	diffLabel := "Custom"
	if g.difficulty != "" {
		diffLabel = string(g.difficulty)
	}
	msg := fmt.Sprintf("Difficulty: %s  Memes: %d  Flags: %d  Status: %s  (LMB reveal, RMB flag, 1/2/3 difficulty, R restart, Esc quit)", diffLabel, g.cfg.MemeCount, flags, status)
	if g.board.Status == board.StatusWon {
		msg = "You win! Press R to restart or Esc to quit."
	} else if g.board.Status == board.StatusLost {
		msg = "You hit a meme. Press R to restart or Esc to quit."
	}
	ebitenutil.DebugPrintAt(screen, msg, 8, 8)
}

func (g *Game) Layout(int, int) (int, int) {
	return g.cfg.Cols * g.cfg.TileSize, g.cfg.Rows * g.cfg.TileSize
}

func (g *Game) reset() error {
	seed := g.cfg.Seed
	if seed == 0 {
		seed = g.rng.Int63()
	}
	b, err := board.New(board.Config{
		Rows:      g.cfg.Rows,
		Cols:      g.cfg.Cols,
		MemeCount: g.cfg.MemeCount,
		Seed:      seed,
	})
	if err != nil {
		return err
	}
	g.board = b
	return nil
}

func (g *Game) applyPreset(preset board.Preset) error {
	if preset == "" {
		return nil
	}
	cfg, err := board.PresetConfig(preset, g.cfg.Seed)
	if err != nil {
		return err
	}
	g.cfg.Rows = cfg.Rows
	g.cfg.Cols = cfg.Cols
	g.cfg.MemeCount = cfg.MemeCount
	g.cfg.Seed = cfg.Seed
	g.difficulty = preset

	ebiten.SetWindowSize(g.cfg.Cols*g.cfg.TileSize, g.cfg.Rows*g.cfg.TileSize)
	return nil
}

func (g *Game) handleDifficultyShortcuts() bool {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		return g.setDifficulty(board.PresetEasy)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		return g.setDifficulty(board.PresetMedium)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		return g.setDifficulty(board.PresetHard)
	}
	return false
}

func (g *Game) setDifficulty(preset board.Preset) bool {
	if preset == "" || preset == g.difficulty {
		return false
	}
	if err := g.applyPreset(preset); err != nil {
		return false
	}
	if err := g.reset(); err != nil {
		return false
	}
	return true
}

func (g *Game) handleReveal() error {
	if g.board.Status != board.StatusActive {
		return nil
	}
	row, col, ok := g.cursorCell()
	if !ok {
		return nil
	}
	_, err := g.board.Reveal(row, col)
	return err
}

func (g *Game) handleFlag() {
	if g.board.Status != board.StatusActive {
		return
	}
	row, col, ok := g.cursorCell()
	if !ok {
		return
	}
	cell := &g.board.Cells[row][col]
	if cell.Revealed {
		return
	}
	cell.Flagged = !cell.Flagged
}

func (g *Game) cursorCell() (int, int, bool) {
	x, y := ebiten.CursorPosition()
	if x < 0 || y < 0 {
		return 0, 0, false
	}
	row := y / g.cfg.TileSize
	col := x / g.cfg.TileSize
	if row < 0 || row >= g.cfg.Rows || col < 0 || col >= g.cfg.Cols {
		return 0, 0, false
	}
	return row, col, true
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	width := g.cfg.Cols * g.cfg.TileSize
	height := g.cfg.Rows * g.cfg.TileSize
	for x := 0; x < width; x += g.cfg.TileSize {
		ebitenutil.DrawRect(screen, float64(x), 0, 1, float64(height), gridColor)
	}
	for y := 0; y < height; y += g.cfg.TileSize {
		ebitenutil.DrawRect(screen, 0, float64(y), float64(width), 1, gridColor)
	}
}

func (g *Game) drawCell(screen *ebiten.Image, row, col int) {
	cell := g.board.Cells[row][col]
	x := float64(col * g.cfg.TileSize)
	y := float64(row * g.cfg.TileSize)

	if cell.Revealed {
		ebitenutil.DrawRect(screen, x, y, float64(g.cfg.TileSize), float64(g.cfg.TileSize), revealedColor)
	} else {
		ebitenutil.DrawRect(screen, x, y, float64(g.cfg.TileSize), float64(g.cfg.TileSize), hiddenColor)
	}

	showMeme := cell.HasMeme && (cell.Revealed || g.board.Status == board.StatusLost)
	if showMeme {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(g.cfg.TileSize)/float64(g.assets.Meme.Bounds().Dx()), float64(g.cfg.TileSize)/float64(g.assets.Meme.Bounds().Dy()))
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.assets.Meme, op)
		return
	}

	if cell.Flagged && !cell.Revealed {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(g.cfg.TileSize)/float64(g.assets.Flag.Bounds().Dx()), float64(g.cfg.TileSize)/float64(g.assets.Flag.Bounds().Dy()))
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.assets.Flag, op)
		return
	}

	if cell.Revealed && cell.AdjacentMemes > 0 {
		label := fmt.Sprintf("%d", cell.AdjacentMemes)
		ebitenutil.DebugPrintAt(screen, label, int(x)+g.cfg.TileSize/2-4, int(y)+g.cfg.TileSize/2-4)
	}
}

func countFlags(b *board.Board) int {
	total := 0
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if b.Cells[r][c].Flagged {
				total++
			}
		}
	}
	return total
}

func statusLabel(status board.Status) string {
	switch status {
	case board.StatusWon:
		return "Won"
	case board.StatusLost:
		return "Lost"
	default:
		return "Active"
	}
}
