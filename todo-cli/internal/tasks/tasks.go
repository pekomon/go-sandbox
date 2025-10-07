package tasks

import "errors"

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID   int
	Text string
	Done bool
}

func Add(taskList []Task, text string) []Task {
	return nil
}

func MarkDone(taskList []Task, id int) ([]Task, error) {
	return nil, ErrTaskNotFound
}

func Remove(taskList []Task, id int) ([]Task, error) {
	return nil, ErrTaskNotFound
}

func Sort(taskList []Task, reverse bool) []Task {
	return nil
}
