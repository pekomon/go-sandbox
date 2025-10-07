package storage_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pekomon/go-sandbox/todo-cli/internal/storage"
	"github.com/pekomon/go-sandbox/todo-cli/internal/tasks"
)

func TestSaveAndLoadTasks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	want := []tasks.Task{
		{ID: 1, Text: "write tests"},
		{ID: 2, Text: "implement features", Done: true},
	}

	if err := storage.SaveTasks(path, want); err != nil {
		t.Fatalf("SaveTasks returned error: %v", err)
	}

	got, err := storage.LoadTasks(path)
	if err != nil {
		t.Fatalf("LoadTasks returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded tasks mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestLoadTasksWithCorruptedJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	if _, err := storage.LoadTasks(path); err == nil {
		t.Fatalf("expected error when loading corrupted JSON")
	}
}

func TestLoadTasksFromEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	tasks, err := storage.LoadTasks(path)
	if err != nil {
		t.Fatalf("LoadTasks returned error: %v", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected no tasks, got %d", len(tasks))
	}
}
