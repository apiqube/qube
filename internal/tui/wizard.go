package tui

// InitAnswers holds the user's responses from the interactive init wizard.
type InitAnswers struct {
	TargetURL    string
	HasSwagger   bool
	SwaggerURL   string
	Plugins      []string
	UseDockerSvc bool
}

// RunInitWizard launches an interactive wizard and returns the user's choices.
// Returns nil if the user cancels.
func RunInitWizard() (*InitAnswers, error) {
	// TODO: implementation using charmbracelet/huh or similar
	//
	// Questions:
	// 1. What's your default target URL? (text input)
	// 2. Do you have an OpenAPI/Swagger spec? (yes/no)
	// 3. If yes: URL or path to spec? (text input with validation)
	// 4. Which protocols to enable? (multi-select: http, grpc, graphql, ws, sql, kafka)
	// 5. Enable docker service management? (yes/no)
	//
	// Return assembled InitAnswers
	return nil, nil
}
