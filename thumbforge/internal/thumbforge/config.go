package thumbforge

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ErrNotImplemented is returned by stubbed ThumbForge functions in the tests-first PR.
var ErrNotImplemented = errors.New("thumbforge: not implemented")

// Size represents a target thumbnail size in pixels.
type Size struct {
	Width  int
	Height int
}

// Config defines the batch thumbnail generation settings.
type Config struct {
	InputDir  string
	OutputDir string
	Size      Size
	Format    string
	Crop      bool
}

// Result reports summary data from a batch run.
type Result struct {
	Count int
}

// ParseSize parses a WxH size string (e.g. 320x240).
func ParseSize(input string) (Size, error) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(input)), "x")
	if len(parts) != 2 {
		return Size{}, fmt.Errorf("thumbforge: invalid size %q", input)
	}
	width, err := strconv.Atoi(parts[0])
	if err != nil || width <= 0 {
		return Size{}, fmt.Errorf("thumbforge: invalid size %q", input)
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil || height <= 0 {
		return Size{}, fmt.Errorf("thumbforge: invalid size %q", input)
	}
	return Size{Width: width, Height: height}, nil
}

// Generate produces thumbnails based on the supplied configuration.
func Generate(cfg Config) (Result, error) {
	if cfg.InputDir == "" {
		return Result{}, fmt.Errorf("thumbforge: input directory required")
	}
	if cfg.OutputDir == "" {
		return Result{}, fmt.Errorf("thumbforge: output directory required")
	}
	if cfg.Size.Width <= 0 || cfg.Size.Height <= 0 {
		return Result{}, fmt.Errorf("thumbforge: invalid size")
	}

	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	if format == "" {
		format = "png"
	}
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return Result{}, err
	}

	entries, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return Result{}, err
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		inPath := filepath.Join(cfg.InputDir, name)
		outPath := filepath.Join(cfg.OutputDir, outputName(name, format))

		if err := generateOne(inPath, outPath, cfg.Size, format); err != nil {
			return Result{}, err
		}
		count++
	}

	return Result{Count: count}, nil
}

func outputName(inputName, format string) string {
	ext := strings.ToLower(filepath.Ext(inputName))
	base := strings.TrimSuffix(inputName, ext)
	outExt := format
	if format == "jpeg" {
		outExt = "jpg"
	}
	return fmt.Sprintf("%s.%s", base, outExt)
}

func generateOne(inputPath, outputPath string, size Size, format string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	src, _, err := image.Decode(inFile)
	if err != nil {
		return err
	}

	dst := resizeNearest(src, size)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "jpg", "jpeg":
		return jpeg.Encode(outFile, dst, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(outFile, dst)
	default:
		return fmt.Errorf("thumbforge: unsupported format %q", format)
	}
}

func resizeNearest(src image.Image, size Size) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for y := 0; y < size.Height; y++ {
		srcY := srcBounds.Min.Y + (y*srcH)/size.Height
		for x := 0; x < size.Width; x++ {
			srcX := srcBounds.Min.X + (x*srcW)/size.Width
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}
