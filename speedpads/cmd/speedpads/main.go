package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/pekomon/go-sandbox/speedpads/internal/game"
)

const sampleRate = 44100

var backgroundColor = color.RGBA{0x10, 0x12, 0x18, 0xff}

var padColors = []color.RGBA{
	{0xf0, 0x6c, 0x4e, 0xff},
	{0xf0, 0xc1, 0x4b, 0xff},
	{0x43, 0xbf, 0x7a, 0xff},
	{0x4a, 0xa8, 0xe8, 0xff},
}

var dimPadColors = []color.RGBA{
	{0x68, 0x2d, 0x21, 0xff},
	{0x6a, 0x58, 0x24, 0xff},
	{0x1d, 0x58, 0x3a, 0xff},
	{0x22, 0x4f, 0x70, 0xff},
}

type app struct {
	state       *game.State
	audio       *audio.Context
	padSounds   []*audio.Player
	gameOver    *audio.Player
	lastLitPad  int
	lastPhase   game.Phase
}

type soundSpec struct {
	frequency float64
	duration  float64
	volume    float64
}

func newApp() (*app, error) {
	st, err := game.New(game.DefaultConfig())
	if err != nil {
		return nil, err
	}

	ctx := audio.NewContext(sampleRate)
	padSounds, err := newPadPlayers(ctx)
	if err != nil {
		return nil, err
	}
	gameOver, err := newGameOverPlayer(ctx)
	if err != nil {
		return nil, err
	}

	return &app{
		state:      st,
		audio:      ctx,
		padSounds:  padSounds,
		gameOver:   gameOver,
		lastLitPad: st.LitPad,
		lastPhase:  st.Phase,
	}, nil
}

func (a *app) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	in := game.Input{
		Start:   inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter),
		Restart: inpututil.IsKeyJustPressed(ebiten.KeyR),
		Pad:     -1,
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyD):
		in.Pad = 0
		in.Press = true
	case inpututil.IsKeyJustPressed(ebiten.KeyF):
		in.Pad = 1
		in.Press = true
	case inpututil.IsKeyJustPressed(ebiten.KeyJ):
		in.Pad = 2
		in.Press = true
	case inpututil.IsKeyJustPressed(ebiten.KeyK):
		in.Pad = 3
		in.Press = true
	}

	if err := a.state.Step(in); err != nil {
		return err
	}

	a.playAudioEvents()
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	padWidth := 120.0
	padHeight := 120.0
	gap := 28.0
	rowWidth := float64(a.state.PadCount)*padWidth + float64(a.state.PadCount-1)*gap
	startX := (a.state.Width - rowWidth) / 2
	startY := 120.0

	for i := 0; i < a.state.PadCount; i++ {
		x := startX + float64(i)*(padWidth+gap)
		y := startY
		padColor := dimPadColors[i]
		if a.state.LitPad == i {
			padColor = padColors[i]
		}
		ebitenutil.DrawRect(screen, x, y, padWidth, padHeight, padColor)
		ebitenutil.DebugPrintAt(screen, padLabel(i), int(x+52), int(y+50))
	}

	header := fmt.Sprintf("Round %d  Score %d  Sequence %d", a.state.Round, a.state.Score, len(a.state.Sequence))
	ebitenutil.DebugPrintAt(screen, header, 16, 14)
	ebitenutil.DebugPrintAt(screen, statusLine(a.state), 16, 34)
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

func newPadPlayers(ctx *audio.Context) ([]*audio.Player, error) {
	specs := []soundSpec{
		{frequency: 392, duration: 0.08, volume: 0.18},
		{frequency: 493.88, duration: 0.08, volume: 0.18},
		{frequency: 587.33, duration: 0.08, volume: 0.18},
		{frequency: 783.99, duration: 0.08, volume: 0.18},
	}

	players := make([]*audio.Player, 0, len(specs))
	for _, spec := range specs {
		player, err := newTonePlayer(ctx, []soundSpec{spec})
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func newGameOverPlayer(ctx *audio.Context) (*audio.Player, error) {
	return newTonePlayer(ctx, []soundSpec{
		{frequency: 220, duration: 0.12, volume: 0.20},
		{frequency: 174.61, duration: 0.12, volume: 0.18},
		{frequency: 130.81, duration: 0.18, volume: 0.18},
	})
}

func newTonePlayer(ctx *audio.Context, specs []soundSpec) (*audio.Player, error) {
	return ctx.NewPlayerFromBytes(synthesizeTone(specs)), nil
}

func synthesizeTone(specs []soundSpec) []byte {
	totalSamples := 0
	for _, spec := range specs {
		totalSamples += int(float64(sampleRate) * spec.duration)
	}

	data := make([]byte, totalSamples*4)
	offset := 0
	for _, spec := range specs {
		segmentSamples := int(float64(sampleRate) * spec.duration)
		for i := 0; i < segmentSamples; i++ {
			progress := float64(i) / float64(segmentSamples)
			envelope := 1.0
			if progress < 0.08 {
				envelope = progress / 0.08
			} else if progress > 0.84 {
				envelope = (1.0 - progress) / 0.16
			}
			if envelope < 0 {
				envelope = 0
			}

			value := math.Sin(2*math.Pi*spec.frequency*float64(i)/sampleRate) * spec.volume * envelope
			sample := int16(value * math.MaxInt16)
			binary.LittleEndian.PutUint16(data[offset:], uint16(sample))
			binary.LittleEndian.PutUint16(data[offset+2:], uint16(sample))
			offset += 4
		}
	}
	return data
}

func (a *app) playAudioEvents() {
	if a.lastPhase == game.PhaseShowing && a.state.Phase == game.PhaseShowing &&
		a.state.LitPad >= 0 && a.state.LitPad != a.lastLitPad {
		a.playPlayer(a.padSounds[a.state.LitPad])
	}

	if a.lastPhase != game.PhaseShowing && a.state.Phase == game.PhaseShowing && a.state.LitPad >= 0 {
		a.playPlayer(a.padSounds[a.state.LitPad])
	}

	if a.lastPhase != game.PhaseGameOver && a.state.Phase == game.PhaseGameOver {
		a.playPlayer(a.gameOver)
	}

	a.lastLitPad = a.state.LitPad
	a.lastPhase = a.state.Phase
}

func (a *app) playPlayer(player *audio.Player) {
	if player == nil {
		return
	}
	player.Pause()
	if err := player.Rewind(); err != nil {
		return
	}
	player.Play()
}

func padLabel(pad int) string {
	switch pad {
	case 0:
		return "D"
	case 1:
		return "F"
	case 2:
		return "J"
	case 3:
		return "K"
	default:
		return "?"
	}
}

func statusLine(st *game.State) string {
	switch st.Phase {
	case game.PhaseAttract:
		return "Press Space or Enter to start. Repeat the sequence on D/F/J/K. Esc quits."
	case game.PhaseShowing:
		return "Watch the pattern. Each round replays the full sequence a bit faster."
	case game.PhaseInput:
		return fmt.Sprintf("Repeat the sequence now. Timeout in %d ticks.", st.InputTimeoutTicks-st.WatchdogTick)
	case game.PhaseGameOver:
		return "Game over. Press R, Space, or Enter to restart."
	default:
		return ""
	}
}
