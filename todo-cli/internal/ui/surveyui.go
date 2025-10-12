package ui

import "github.com/AlecAivazis/survey/v2"

type SurveyUI struct{}

func (SurveyUI) Select(title string, options []string) (int, error) {
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
