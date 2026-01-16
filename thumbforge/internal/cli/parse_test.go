package cli_test

import (
	"os"
	"testing"

	"github.com/pekomon/go-sandbox/thumbforge/internal/cli"
	"github.com/pekomon/go-sandbox/thumbforge/internal/thumbforge"
)

func TestParseArgs(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{
		"--in", inDir,
		"--out", outDir,
		"--size", "120x90",
		"--format", "jpg",
		"--crop",
	}

	cfg, err := cli.ParseArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := thumbforge.Config{
		InputDir:  inDir,
		OutputDir: outDir,
		Size:      thumbforge.Size{Width: 120, Height: 90},
		Format:    "jpg",
		Crop:      true,
	}

	if cfg != want {
		t.Fatalf("unexpected config: got %+v want %+v", cfg, want)
	}
}

func TestParseArgsWithWidthHeight(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{
		"--in", inDir,
		"--out", outDir,
		"--width", "80",
		"--height", "60",
		"--format", "png",
	}

	cfg, err := cli.ParseArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := thumbforge.Config{
		InputDir:  inDir,
		OutputDir: outDir,
		Size:      thumbforge.Size{Width: 80, Height: 60},
		Format:    "png",
	}

	if cfg != want {
		t.Fatalf("unexpected config: got %+v want %+v", cfg, want)
	}
}

func TestParseArgsMissingRequired(t *testing.T) {
	args := []string{"--size", "120x90"}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestParseArgsRejectsMixedSizeFlags(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{
		"--in", inDir,
		"--out", outDir,
		"--size", "120x90",
		"--width", "120",
		"--height", "90",
	}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestParseArgsInvalidSize(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{"--in", inDir, "--out", outDir, "--size", "abc"}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestParseArgsInvalidFormat(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{
		"--in", inDir,
		"--out", outDir,
		"--size", "120x90",
		"--format", "gif",
	}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestParseArgsDefaults(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{"--in", inDir, "--out", outDir, "--size", "120x90"}

	cfg, err := cli.ParseArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Format == "" {
		t.Fatalf("expected default format")
	}
	if cfg.InputDir != inDir || cfg.OutputDir != outDir {
		t.Fatalf("unexpected directories")
	}
	if cfg.Size != (thumbforge.Size{Width: 120, Height: 90}) {
		t.Fatalf("unexpected size")
	}
}

func TestParseArgsMissingWidthHeight(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()

	args := []string{
		"--in", inDir,
		"--out", outDir,
		"--width", "120",
	}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestParseArgsRejectsMissingInputDir(t *testing.T) {
	args := []string{"--out", os.TempDir(), "--size", "120x90"}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Fatalf("expected error")
	}
}
