package output

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/ui"
)

// Pretty is an EventHandler that prints structured log lines as engine
// events arrive. Used by qube run when stdout isn't a TTY (CI, pipes) and
// by --no-tui. Each event maps to one log line; failures get extra
// indented lines for assertion details.
type Pretty struct {
	w       io.Writer
	mu      sync.Mutex
	verbose bool
	logger  *ui.Logger
}

// NewPretty creates a Pretty output handler writing to w. If verbose, it
// emits TestStarted lines (else tests show up only on completion).
func NewPretty(w io.Writer, verbose bool) *Pretty {
	logger := ui.DefaultLogger.WithWriter(w)
	logger.SetVerbose(verbose)
	return &Pretty{w: w, verbose: verbose, logger: logger}
}

// Handle dispatches one event to the right log writer.
func (p *Pretty) Handle(event engine.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch e := event.(type) {
	case engine.RunStarted:
		p.logger.Logf(ui.LevelInfo, "qube run · %s · %s · %s",
			plural(e.TotalTests, "test"),
			plural(len(e.Files), "file"),
			plural(e.TotalWaves, "wave"))
	case engine.WaveStarted:
		mode := "sequential"
		if e.Parallel {
			mode = "parallel"
		}
		p.logger.Logf(ui.LevelDebug, "wave %d (%s)", e.Index+1, mode)
	case engine.TestStarted:
		if p.verbose {
			p.logger.Logf(ui.LevelDebug, "▶ %s", e.Name)
		}
	case engine.TestCompleted:
		p.onTestCompleted(e)
	case engine.RunCompleted:
		p.onRunCompleted(e)
	}
}

func (p *Pretty) onTestCompleted(e engine.TestCompleted) {
	dur := formatDuration(e.Duration)
	switch e.Status {
	case engine.StatusPassed:
		p.logger.Logf(ui.LevelSuccess, "✓ %s (%s)", e.Name, dur)
	case engine.StatusSkipped:
		p.logger.Logf(ui.LevelInfo, "⏭ %s (skipped)", e.Name)
	case engine.StatusErrored:
		p.logger.Logf(ui.LevelError, "✗ %s (%s)\n%s", e.Name, dur, errorDetails(e))
	case engine.StatusFailed:
		p.logger.Logf(ui.LevelError, "✗ %s (%s)\n%s", e.Name, dur, assertionDetails(e))
	}
}

func (p *Pretty) onRunCompleted(e engine.RunCompleted) {
	p.logger.Logf(ui.LevelInfo, "Run finished · %s · %s",
		summaryLine(e),
		formatDuration(e.Duration))
}

// summaryLine returns "3 passed, 1 failed, 0 skipped" with sensible coloring.
func summaryLine(e engine.RunCompleted) string {
	parts := []string{ui.SuccessStyle.Render(fmt.Sprintf("%d passed", e.Passed))}
	if e.Failed > 0 {
		parts = append(parts, ui.FailureStyle.Render(fmt.Sprintf("%d failed", e.Failed)))
	}
	if e.Errored > 0 {
		parts = append(parts, ui.FailureStyle.Render(fmt.Sprintf("%d errored", e.Errored)))
	}
	if e.Skipped > 0 {
		parts = append(parts, ui.MutedStyle.Render(fmt.Sprintf("%d skipped", e.Skipped)))
	}
	return joinWithSep(parts, ", ")
}

func joinWithSep(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

func errorDetails(e engine.TestCompleted) string {
	if e.Error == "" {
		return ui.MutedStyle.Render("  (no error message)")
	}
	return "  " + ui.MutedStyle.Render(e.Error)
}

func assertionDetails(e engine.TestCompleted) string {
	out := ""
	for i, a := range e.Assertions {
		if a.Passed {
			continue
		}
		if i > 0 {
			out += "\n"
		}
		out += "  "
		out += ui.FailureStyle.Render(a.Expression) + "  "
		out += ui.MutedStyle.Render(fmt.Sprintf("expected %v, actual %v", a.Expected, a.Actual))
	}
	if out == "" && e.Error != "" {
		return "  " + ui.MutedStyle.Render(e.Error)
	}
	if out == "" {
		return ui.MutedStyle.Render("  (no failure details)")
	}
	return out
}

func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

func formatDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%dms", d.Milliseconds())
	default:
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
}
