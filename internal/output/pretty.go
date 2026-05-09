package output

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/ui"
)

// Pretty is an EventHandler that prints styled, progressive output as events
// arrive. It's the default for non-interactive shells (CI, pipes) and the
// fallback when the user explicitly disables the live TUI.
type Pretty struct {
	w       io.Writer
	mu      sync.Mutex
	verbose bool

	files int
	tests int
	waves int
	start time.Time
}

// NewPretty creates a Pretty output handler writing to w.
func NewPretty(w io.Writer, verbose bool) *Pretty {
	return &Pretty{w: w, verbose: verbose}
}

// Handle dispatches one event to the right renderer.
func (p *Pretty) Handle(event engine.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch e := event.(type) {
	case engine.RunStarted:
		p.onRunStarted(e)
	case engine.WaveStarted:
		p.onWaveStarted(e)
	case engine.TestStarted:
		if p.verbose {
			p.onTestStarted(e)
		}
	case engine.TestCompleted:
		p.onTestCompleted(e)
	case engine.RunCompleted:
		p.onRunCompleted(e)
	}
}

func (p *Pretty) onRunStarted(e engine.RunStarted) {
	p.files = len(e.Files)
	p.tests = e.TotalTests
	p.waves = e.TotalWaves
	p.start = time.Now()

	header := ui.Brand.Render("qube run")
	subtitle := ui.Muted.Render(fmt.Sprintf("%d tests · %d files · %d waves",
		e.TotalTests, len(e.Files), e.TotalWaves))
	fmt.Fprintln(p.w, ui.Header.Render(header+"  "+subtitle))
}

func (p *Pretty) onWaveStarted(e engine.WaveStarted) {
	if !p.verbose {
		return
	}
	label := fmt.Sprintf("Wave %d", e.Index+1)
	if e.Parallel {
		label += "  · parallel"
	}
	fmt.Fprintln(p.w, ui.Accent.Render(label))
}

func (p *Pretty) onTestStarted(e engine.TestStarted) {
	icon := ui.StatusRunning.Icon()
	line := fmt.Sprintf("%s  %s  %s",
		ui.Accent.Render(icon),
		e.Name,
		ui.Muted.Render(string(e.Protocol)+" "+e.Target),
	)
	fmt.Fprintln(p.w, line)
}

func (p *Pretty) onTestCompleted(e engine.TestCompleted) {
	status := mapStatus(e.Status)
	icon := ui.StatusStyle(status).Render(status.Icon())
	durStr := ui.Muted.Render(formatDuration(e.Duration))
	line := fmt.Sprintf("%s  %s  %s", icon, e.Name, durStr)
	fmt.Fprintln(p.w, line)

	if e.Status == engine.StatusFailed || e.Status == engine.StatusErrored {
		p.printFailureDetails(e)
	}
}

func (p *Pretty) printFailureDetails(e engine.TestCompleted) {
	if e.Error != "" {
		fmt.Fprintf(p.w, "    %s %s\n",
			ui.Failure.Render("error:"),
			e.Error,
		)
	}
	for _, a := range e.Assertions {
		if a.Passed {
			continue
		}
		exp := ui.Muted.Render(fmt.Sprintf("expected: %v", a.Expected))
		act := ui.Muted.Render(fmt.Sprintf("actual:   %v", a.Actual))
		fmt.Fprintf(p.w, "    %s  %s\n", ui.Failure.Render("✗"), a.Expression)
		fmt.Fprintf(p.w, "      %s\n", exp)
		fmt.Fprintf(p.w, "      %s\n", act)
		if a.Message != "" {
			fmt.Fprintf(p.w, "      %s\n", ui.Muted.Render(a.Message))
		}
	}
}

func (p *Pretty) onRunCompleted(e engine.RunCompleted) {
	body := strings.Builder{}
	body.WriteString(ui.Success.Render(fmt.Sprintf("%d passed", e.Passed)))
	if e.Failed > 0 {
		body.WriteString("   ")
		body.WriteString(ui.Failure.Render(fmt.Sprintf("%d failed", e.Failed)))
	}
	if e.Errored > 0 {
		body.WriteString("   ")
		body.WriteString(ui.Failure.Render(fmt.Sprintf("%d errored", e.Errored)))
	}
	if e.Skipped > 0 {
		body.WriteString("   ")
		body.WriteString(ui.Muted.Render(fmt.Sprintf("%d skipped", e.Skipped)))
	}
	body.WriteString("\n")
	body.WriteString(ui.Muted.Render(formatDuration(e.Duration) + " total"))
	fmt.Fprintln(p.w)
	fmt.Fprintln(p.w, ui.SummaryCard.Render(body.String()))
}

func mapStatus(s engine.TestStatus) ui.Status {
	switch s {
	case engine.StatusPassed:
		return ui.StatusPassed
	case engine.StatusFailed:
		return ui.StatusFailed
	case engine.StatusSkipped:
		return ui.StatusSkipped
	case engine.StatusErrored:
		return ui.StatusErrored
	}
	return ui.StatusPending
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
