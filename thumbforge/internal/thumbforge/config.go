package thumbforge

import "errors"

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
	return Size{}, ErrNotImplemented
}

// Generate produces thumbnails based on the supplied configuration.
func Generate(cfg Config) (Result, error) {
	return Result{}, ErrNotImplemented
}
