package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/pekomon/go-sandbox/todo-cli/internal/ui"
)

const (
	menuIndexList = iota
	menuIndexAdd
	menuIndexDone
	menuIndexRemove
	menuIndexClear
	menuIndexExit
)

func TestMenuSelectionsViaUI(t *testing.T) {
	t.Setenv("TODO_CLI_PATH", filepath.Join(t.TempDir(), "tasks.json"))

	fake := &ui.Fake{
		Choices: []int{
			menuIndexList,
			menuIndexAdd,
			menuIndexDone,
			menuIndexRemove,
			menuIndexClear,
			menuIndexExit,
		},
	}

	stdout, stderr, exit := runMenuWithFakeUI(t, fake)

	if exit != 0 {
		t.Fatalf("expected exit code 0, got %d. stdout=%q stderr=%q", exit, stdout, stderr)
	}

	wantStdout := []string{
		"No tasks found.",
		"added #1",
		"done #1",
		"removed #1",
		"cleared",
		"Goodbye!",
	}

	for _, want := range wantStdout {
		if !strings.Contains(stdout, want) {
			t.Errorf("expected stdout to contain %q, got:\n%s", want, stdout)
		}
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr output:\n%s", stderr)
	}
}

func runMenuWithFakeUI(t *testing.T, fake ui.MenuUI) (string, string, int) {
	t.Helper()

	previous := menuUI
	menuUI = fake
	defer func() {
		menuUI = previous
	}()

	return runMenuHarness(t, []string{"menu"}, "")
}
