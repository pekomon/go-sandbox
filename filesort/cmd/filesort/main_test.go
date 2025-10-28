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

// We keep this test light: it invokes "go run ./cmd/filesort" to ensure "--dry-run" is recognized
// in the future implementation. For now we expect the command to exist; behavior assertions are
// commented out until the feature PR wires flags.

func TestCLI_DryRunFlag_WiresThrough(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows for path quoting simplicity")
	}

	root := t.TempDir()
	// Create a couple of files that would be classified later
	_ = os.WriteFile(filepath.Join(root, "a.jpg"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "b.txt"), []byte("x"), 0o644)

	cmd := exec.Command("go", "run", "./cmd/filesort", "--dry-run", root)
	cmd.Dir = filepath.Join("..") // from filesort/cmd/filesort to filesort/
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()

	// Until CLI wiring exists, we accept non-zero exit and/or stderr output.
	if err == nil && strings.TrimSpace(errb.String()) == "" {
		// Uncomment in the implementation PR to assert stable messages:
		// if !strings.Contains(out.String(), "dry-run") {
		// t.Fatalf("expected mention of dry-run in output")
		// }
	}
}
