package ui

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// Assets holds the UI sprites loaded from disk.
type Assets struct {
	Meme *ebiten.Image
	Flag *ebiten.Image
}

// LoadAssets locates and loads the required placeholder sprites.
func LoadAssets() (*Assets, error) {
	meme, err := loadImage("meme.png")
	if err != nil {
		return nil, err
	}
	flag, err := loadImage("flag.png")
	if err != nil {
		return nil, err
	}
	return &Assets{Meme: meme, Flag: flag}, nil
}

func loadImage(name string) (*ebiten.Image, error) {
	path, err := resolveAssetPath(name)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ui: open asset %s: %w", name, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("ui: decode asset %s: %w", name, err)
	}
	return ebiten.NewImageFromImage(img), nil
}

func resolveAssetPath(name string) (string, error) {
	candidates := []string{
		filepath.Join("assets", name),
	}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "assets", name),
			filepath.Join(exeDir, "..", "assets", name),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("ui: asset %s not found", name)
}
