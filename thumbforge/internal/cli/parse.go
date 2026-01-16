package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/pekomon/go-sandbox/thumbforge/internal/thumbforge"
)

// ParseArgs parses CLI arguments into a ThumbForge configuration.
func ParseArgs(args []string) (thumbforge.Config, error) {
	fs := flag.NewFlagSet("thumbforge", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inputDir string
	var outputDir string
	var sizeRaw string
	var width int
	var height int
	var format string
	var crop bool

	fs.StringVar(&inputDir, "in", "", "input directory")
	fs.StringVar(&outputDir, "out", "", "output directory")
	fs.StringVar(&sizeRaw, "size", "", "thumbnail size (WxH)")
	fs.IntVar(&width, "width", 0, "thumbnail width in pixels")
	fs.IntVar(&height, "height", 0, "thumbnail height in pixels")
	fs.StringVar(&format, "format", "png", "output format (png, jpg)")
	fs.BoolVar(&crop, "crop", false, "center-crop before resizing")

	if err := fs.Parse(args); err != nil {
		return thumbforge.Config{}, err
	}
	if inputDir == "" {
		return thumbforge.Config{}, fmt.Errorf("thumbforge: input directory required")
	}
	if outputDir == "" {
		return thumbforge.Config{}, fmt.Errorf("thumbforge: output directory required")
	}
	if sizeRaw != "" && (width > 0 || height > 0) {
		return thumbforge.Config{}, fmt.Errorf("thumbforge: size and width/height are mutually exclusive")
	}

	var size thumbforge.Size
	if sizeRaw != "" {
		parsed, err := thumbforge.ParseSize(sizeRaw)
		if err != nil {
			return thumbforge.Config{}, err
		}
		size = parsed
	} else {
		if width <= 0 || height <= 0 {
			return thumbforge.Config{}, fmt.Errorf("thumbforge: size required")
		}
		size = thumbforge.Size{Width: width, Height: height}
	}

	normalizedFormat, err := thumbforge.NormalizeFormat(format)
	if err != nil {
		return thumbforge.Config{}, err
	}

	return thumbforge.Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Size:      size,
		Format:    normalizedFormat,
		Crop:      crop,
	}, nil
}
