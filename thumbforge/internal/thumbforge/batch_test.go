package thumbforge_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/pekomon/go-sandbox/thumbforge/internal/thumbforge"
)

func TestGenerateBatch(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	if err := writePNG(filepath.Join(inDir, "one.png"), 4, 4); err != nil {
		t.Fatalf("write png: %v", err)
	}
	if err := writePNG(filepath.Join(inDir, "two.png"), 8, 8); err != nil {
		t.Fatalf("write png: %v", err)
	}

	cfg := thumbforge.Config{
		InputDir:  inDir,
		OutputDir: outDir,
		Size:      thumbforge.Size{Width: 2, Height: 2},
		Format:    "png",
	}

	result, err := thumbforge.Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Fatalf("unexpected count: got %d want 2", result.Count)
	}

	if _, err := os.Stat(filepath.Join(outDir, "one.png")); err != nil {
		t.Fatalf("expected thumbnail for one.png: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "two.png")); err != nil {
		t.Fatalf("expected thumbnail for two.png: %v", err)
	}
}

func TestGenerateEmptyInput(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	cfg := thumbforge.Config{
		InputDir:  inDir,
		OutputDir: outDir,
		Size:      thumbforge.Size{Width: 2, Height: 2},
		Format:    "png",
	}

	_, err := thumbforge.Generate(cfg)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func writePNG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 0x33, G: 0x66, B: 0x99, A: 0xff})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
