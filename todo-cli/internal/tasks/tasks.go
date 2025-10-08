package tasks

import (
	"errors"
	"sort"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	// CreatedAt intentionally omitted from JSON schema for now to keep tests simple;
	// newest-first is implemented via ID ordering.
}

var ErrNotFound = errors.New("task not found")

// ErrTaskNotFound is kept for backward compatibility with earlier versions of
// the package. It aliases ErrNotFound so older code (and tests) continue to
// compile while newer code can use the shorter name.
var ErrTaskNotFound = ErrNotFound

// NextID returns the next sequential ID (max(existing)+1), starting at 1.
func NextID(list []Task) int {
	max := 0
	for _, t := range list {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

// Add appends a new task with a new ID.
func Add(list []Task, text string) []Task {
	id := NextID(list)
	return append(list, Task{ID: id, Text: text, Done: false})
}

// MarkDone sets Done=true for the given id.
func MarkDone(list []Task, id int) ([]Task, error) {
	for i := range list {
		if list[i].ID == id {
			list[i].Done = true
			return list, nil
		}
	}
	return list, ErrNotFound
}

// Remove deletes the task with the given id.
func Remove(list []Task, id int) ([]Task, error) {
	for i := range list {
		if list[i].ID == id {
			return append(list[:i], list[i+1:]...), nil
		}
	}
	return list, ErrNotFound
}

// Clear removes all tasks.
func Clear(_ []Task) []Task { return nil }

// SortNewestFirst returns a new slice sorted by ID desc (newest first).
func SortNewestFirst(list []Task) []Task {
	cp := append([]Task(nil), list...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].ID > cp[j].ID })
	return cp
}

// SortOldestFirst returns a new slice sorted by ID asc.
func SortOldestFirst(list []Task) []Task {
	cp := append([]Task(nil), list...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].ID < cp[j].ID })
	return cp
}

// Sort sorts tasks newest-first by default, or oldest-first when reverse is
// true. It preserves the original slice by returning a copy.
func Sort(list []Task, reverse bool) []Task {
	if reverse {
		return SortOldestFirst(list)
	}
	return SortNewestFirst(list)
}
