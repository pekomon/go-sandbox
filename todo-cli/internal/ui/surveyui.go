package ui

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

var ErrNonInteractive = errors.New("survey ui requires an interactive terminal")

type SurveyUI struct{}

func (SurveyUI) Select(title string, options []string) (int, error) {
	if !isInteractiveTerminal(os.Stdin) || !isInteractiveTerminal(os.Stdout) {
		return 0, ErrNonInteractive
	}

	var choice string
	prompt := &survey.Select{
		Message:  title,
		Options:  options,
		PageSize: 10,
	}
	if err := survey.AskOne(prompt, &choice, survey.WithValidator(survey.Required)); err != nil {
		return 0, err
	}
	for i, s := range options {
		if s == choice {
			return i, nil
		}
	}
	return 0, nil
}

func isInteractiveTerminal(file *os.File) bool {
	if file == nil {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
