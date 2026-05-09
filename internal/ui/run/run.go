package run

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/apiqube/engine"
)

// Run executes eng.Run inside a Bubble Tea program. The TUI receives engine
// events via a tea.Program.Send bridge and renders a live view; when the
// engine finishes the model receives engineDoneMsg and the program exits.
//
// Returns the engine.Results produced by Run and the engine-level error
// (parse failure, ctx cancel, etc.). The TUI itself does not return its own
// error — only catastrophic Bubble Tea failures.
func Run(ctx context.Context, eng *engine.Engine, input engine.Input, opts ...engine.RunOption) (*engine.Results, error) {
	model := New()

	prog := tea.NewProgram(
		model,
		tea.WithContext(ctx),
	)

	bridge := &handler{program: prog}

	var (
		results *engine.Results
		runErr  error
		wg      sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runOpts := append([]engine.RunOption{engine.WithHandler(bridge)}, opts...)
		results, runErr = eng.Run(ctx, input, runOpts...)
		prog.Send(engineDoneMsg{results: results, err: runErr})
	}()

	if _, err := prog.Run(); err != nil {
		// Tea failed (e.g., not a TTY). Drain the engine to avoid a leak.
		wg.Wait()
		if runErr != nil {
			return results, runErr
		}
		return results, err
	}
	wg.Wait()

	return results, runErr
}
