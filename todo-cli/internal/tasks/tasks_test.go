package tasks_test

import (
	"errors"
	"testing"

	"github.com/pekomon/go-sandbox/todo-cli/internal/tasks"
)

func TestAddAssignsSequentialIDs(t *testing.T) {
	list := []tasks.Task{}

	list = tasks.Add(list, "write tests")
	list = tasks.Add(list, "implement features")
	list = tasks.Add(list, "ship it")

	if len(list) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(list))
	}

	for i, task := range list {
		wantID := i + 1
		if task.ID != wantID {
			t.Fatalf("task %q expected ID %d, got %d", task.Text, wantID, task.ID)
		}
	}
}

func TestMarkDoneUpdatesStatus(t *testing.T) {
	list := []tasks.Task{}
	list = tasks.Add(list, "write tests")
	list = tasks.Add(list, "implement features")

	updated, err := tasks.MarkDone(list, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !updated[1].Done {
		t.Fatalf("expected task 2 to be marked done")
	}

	if updated[0].Done {
		t.Fatalf("expected task 1 to remain not done")
	}

	_, err = tasks.MarkDone(list, 99)
	if !errors.Is(err, tasks.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestRemoveDeletesTask(t *testing.T) {
	list := []tasks.Task{}
	list = tasks.Add(list, "write tests")
	list = tasks.Add(list, "implement features")

	updated, err := tasks.Remove(list, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(updated) != 1 {
		t.Fatalf("expected 1 task remaining, got %d", len(updated))
	}

	if updated[0].ID != 2 {
		t.Fatalf("expected remaining task to have ID 2, got %d", updated[0].ID)
	}

	_, err = tasks.Remove(list, 99)
	if !errors.Is(err, tasks.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestSortOrdersByNewestAndReverse(t *testing.T) {
	list := []tasks.Task{
		{ID: 1, Text: "first"},
		{ID: 2, Text: "second"},
		{ID: 3, Text: "third"},
	}

	sorted := tasks.Sort(list, false)
	if sorted[0].ID != 3 || sorted[1].ID != 2 || sorted[2].ID != 1 {
		t.Fatalf("expected newest-first order, got %+v", sorted)
	}

	reversed := tasks.Sort(list, true)
	if reversed[0].ID != 1 || reversed[1].ID != 2 || reversed[2].ID != 3 {
		t.Fatalf("expected ascending order when reverse flag set, got %+v", reversed)
	}
}
