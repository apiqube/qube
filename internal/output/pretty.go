package output

import "github.com/apiqube/engine"

// Pretty is an EventHandler that renders events to the terminal using pterm + lipgloss.
// It shows progress bars, styled tables, spinners, and colored status icons.
type Pretty struct {
	verbose bool
}

// NewPretty creates a new Pretty output handler.
func NewPretty(verbose bool) *Pretty {
	return &Pretty{verbose: verbose}
}

// Handle dispatches an event to the appropriate renderer.
func (p *Pretty) Handle(event engine.Event) {
	// TODO: implementation
	//
	// Type-switch on event:
	//
	// case engine.RunStarted:
	//   Print header with total tests, wave count
	//
	// case engine.WaveStarted:
	//   Start spinner or progress bar for this wave
	//
	// case engine.TestStarted:
	//   If verbose: log "running: <name>"
	//
	// case engine.TestCompleted:
	//   Render icon (✓/✗/⚠) + name + duration
	//   If failed: show assertion failures with expected/actual diff
	//
	// case engine.WaveCompleted:
	//   Update wave progress
	//
	// case engine.RunCompleted:
	//   Print final summary table
	//
	// case engine.PluginLoaded, ConfigLoaded:
	//   Log in debug mode only
}
