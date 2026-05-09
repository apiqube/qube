// Package run hosts the Bubble Tea program that drives the live UI for
// `qube run`. It bridges engine events into tea messages and renders a
// live table, progress bar, and summary card.
package run

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/apiqube/engine"
)

// eventMsg wraps an engine.Event for delivery into the bubbletea Update loop.
type eventMsg struct{ event engine.Event }

// engineDoneMsg signals the engine.Run goroutine has returned. It carries
// the final Results pointer (may be nil) and any engine-level error.
type engineDoneMsg struct {
	results *engine.Results
	err     error
}

// handler implements engine.EventHandler by sending an eventMsg into a
// running tea.Program. Concurrency-safe: tea.Program.Send is internally
// synchronized.
type handler struct {
	program *tea.Program
}

// Handle is invoked from the engine's runner goroutine.
func (h *handler) Handle(event engine.Event) {
	if h.program == nil {
		return
	}
	h.program.Send(eventMsg{event: event})
}
