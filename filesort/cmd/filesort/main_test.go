package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLI_DryRunFlag_WiresThrough(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows for path quoting simplicity")
	}

	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "a.jpg"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "b.txt"), []byte("x"), 0o644)

	// Run from module root so `go run ./cmd/filesort` resolves correctly.
	cmd := exec.Command("go", "run", "./cmd/filesort", "--dry-run", root)
	cmd.Dir = filepath.Join("..") // from filesort/cmd/filesort to filesort/
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		t.Fatalf("go run: %v, stderr: %s", err, errb.String())
	}
	if !strings.Contains(out.String(), "dry-run") {
		t.Fatalf("expected 'dry-run' mention in output; got: %s", out.String())
	}
}
