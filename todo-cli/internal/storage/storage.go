package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pekomon/go-sandbox/todo-cli/internal/tasks"
)

const (
	envPath    = "TODO_CLI_PATH"
	defaultRel = ".todo-cli/tasks.json"
	lockName   = "tasks.lock"
)

type Lock struct {
	path string
}

// DefaultPath returns the file path for tasks.json: $TODO_CLI_PATH or $HOME/.todo-cli/tasks.json.
func DefaultPath() (string, error) {
	if p := os.Getenv(envPath); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultRel), nil
}

// ensureDir ensures that the directory for the given file exists.
func ensureDir(file string) error {
	dir := filepath.Dir(file)
	return os.MkdirAll(dir, 0o755)
}

// AcquireLock tries to create a lockfile next to the JSON file.
// Best-effort: if it already exists, return an error.
func AcquireLock(jsonPath string) (*Lock, error) {
	dir := filepath.Dir(jsonPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	lp := filepath.Join(dir, lockName)
	f, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("another process may be running (lock exists): %w", err)
	}
	// tiny bit of metadata (timestamp) to help debugging, ignore errors
	_, _ = f.WriteString(time.Now().Format(time.RFC3339Nano))
	_ = f.Close()
	return &Lock{path: lp}, nil
}

// Release removes the lock file.
func (l *Lock) Release() {
	if l == nil || l.path == "" {
		return
	}
	_ = os.Remove(l.path)
}

// LoadTasks loads tasks from jsonPath. If the file doesn't exist, returns empty list.
func LoadTasks(jsonPath string) ([]tasks.Task, error) {
	b, err := os.ReadFile(jsonPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	if len(b) == 0 {
		return nil, nil
	}
	var list []tasks.Task
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// SaveTasks writes tasks to jsonPath (pretty JSON).
func SaveTasks(jsonPath string, list []tasks.Task) error {
	if err := ensureDir(jsonPath); err != nil {
		return err
	}
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := jsonPath + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, jsonPath)
}
