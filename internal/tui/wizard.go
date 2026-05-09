// Package tui hosts interactive prompt flows that don't belong to a specific
// command's runtime UI. Today this is the init wizard.
package tui

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
)

// InitAnswers holds the user's responses from the interactive init wizard.
type InitAnswers struct {
	TargetURL  string
	HasSwagger bool
	SwaggerURL string
	Plugins    []string
}

// RunInitWizard launches an interactive wizard and returns the user's choices.
// Returns nil InitAnswers and a wrapped huh error if the user cancels.
func RunInitWizard() (*InitAnswers, error) {
	answers := &InitAnswers{
		TargetURL: "http://localhost:8080",
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Default target URL").
				Description("The base URL most of your tests will hit.").
				Value(&answers.TargetURL).
				Validate(validateURL),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you have an OpenAPI/Swagger spec?").
				Description("If yes, qube will scaffold tests from it (preview).").
				Value(&answers.HasSwagger),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Path or URL to your spec").
				Value(&answers.SwaggerURL).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("required when you have a spec")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return !answers.HasSwagger }),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Which protocols to enable?").
				Description("Only http is fully wired in v1.0; others are previews.").
				Options(
					huh.NewOption("http  (live)", "http").Selected(true),
					huh.NewOption("grpc  (preview)", "grpc"),
					huh.NewOption("graphql  (preview)", "graphql"),
					huh.NewOption("ws  (preview)", "ws"),
					huh.NewOption("sql  (preview)", "sql"),
					huh.NewOption("kafka  (preview)", "kafka"),
				).
				Value(&answers.Plugins),
		),
	)
	if err := form.Run(); err != nil {
		return nil, err
	}
	return answers, nil
}

func validateURL(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("required")
	}
	if !strings.Contains(s, "://") {
		return errors.New("must include a scheme like http:// or https://")
	}
	return nil
}
