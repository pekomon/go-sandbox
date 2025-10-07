package main

import "errors"

var ErrUnknownCommand = errors.New("unknown command")

type Command struct {
	Name string
	Args []string
}

func ParseCommand(args []string) (Command, error) {
	return Command{}, ErrUnknownCommand
}

func main() {}
