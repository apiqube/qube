package run

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/ui"
)

// Model is the Bubble Tea model for `qube run`. The layout is intentionally
// flat — a brand line, a list of test rows, a progress line — no rounded
// borders or boxed cards.
type Model struct {
	files      []string
	totalTests int
	totalWaves int
	startTime  time.Time

	currentWave    int
	completedWaves int

	testOrder []string
	testState map[string]testEntry

	width int

	spinner  spinner.Model
	progress progress.Model

	finished bool
	summary  *engine.RunCompleted
	results  *engine.Results
	runErr   error
}

type testEntry struct {
	name     string
	file     string
	protocol engine.Protocol
	target   string
	status   ui.Status
	duration time.Duration
	err      string
	failures []engine.AssertionResult
}

// New returns a fresh Model ready to be passed to tea.NewProgram.
func New() *Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = ui.AccentStyle

	pb := progress.New(progress.WithGradient("#5BC8E8", "#C084FC"))
	pb.Width = 36

	return &Model{
		testState: map[string]testEntry{},
		spinner:   sp,
		progress:  pb,
		width:     100,
	}
}

// Init starts the spinner.
func (m *Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update integrates one tea.Msg into the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		if m.width > 8 {
			m.progress.Width = m.width / 2
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case eventMsg:
		return m.handleEngineEvent(msg.event)

	case engineDoneMsg:
		m.results = msg.results
		m.runErr = msg.err
		m.finished = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleEngineEvent(e engine.Event) (tea.Model, tea.Cmd) {
	switch ev := e.(type) {
	case engine.RunStarted:
		m.files = append([]string(nil), ev.Files...)
		m.totalTests = ev.TotalTests
		m.totalWaves = ev.TotalWaves
		m.startTime = time.Now()

	case engine.WaveStarted:
		m.currentWave = ev.Index + 1

	case engine.WaveCompleted:
		m.completedWaves++

	case engine.TestStarted:
		key := ev.Name
		if _, ok := m.testState[key]; !ok {
			m.testOrder = append(m.testOrder, key)
		}
		m.testState[key] = testEntry{
			name:     ev.Name,
			file:     ev.File,
			protocol: ev.Protocol,
			target:   ev.Target,
			status:   ui.StatusRunning,
		}

	case engine.TestCompleted:
		key := ev.Name
		entry := m.testState[key]
		entry.name = ev.Name
		entry.duration = ev.Duration
		entry.status = mapStatus(ev.Status)
		entry.err = ev.Error
		entry.failures = failedAssertions(ev.Assertions)
		m.testState[key] = entry

	case engine.RunCompleted:
		c := ev
		m.summary = &c
	}
	return m, nil
}

// View renders the current model state.
func (m *Model) View() string {
	if m.finished {
		return m.viewSummary()
	}
	return m.viewLive()
}

func (m *Model) viewLive() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteByte('\n')
	b.WriteString(m.renderTestList())
	b.WriteByte('\n')
	b.WriteString(m.renderProgress())
	b.WriteByte('\n')
	return b.String()
}

func (m *Model) viewSummary() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteByte('\n')
	if len(m.testOrder) > 0 {
		b.WriteString(m.renderTestList())
		b.WriteByte('\n')
	}
	b.WriteString(m.renderSummary())
	b.WriteByte('\n')
	if details := m.renderFailureDetails(); details != "" {
		b.WriteByte('\n')
		b.WriteString(details)
		b.WriteByte('\n')
	}
	if m.runErr != nil {
		b.WriteString(ui.FailureStyle.Render("engine error: " + m.runErr.Error()))
		b.WriteByte('\n')
	}
	return b.String()
}

// renderHeader is one bold brand line — no border, no padding.
func (m *Model) renderHeader() string {
	title := ui.BrandStyle.Render("qube run")
	subtitle := ui.MutedStyle.Render(fmt.Sprintf(" · %d tests · %d files · %d waves",
		m.totalTests, len(m.files), m.totalWaves))
	return title + subtitle
}

func (m *Model) renderTestList() string {
	if len(m.testOrder) == 0 {
		return ui.MutedStyle.Render("waiting for first test...")
	}
	var lines []string
	for _, name := range m.testOrder {
		t := m.testState[name]
		lines = append(lines, m.renderRow(t))
	}
	return strings.Join(lines, "\n")
}

func (m *Model) renderRow(t testEntry) string {
	icon := t.status.Icon()
	if t.status == ui.StatusRunning {
		icon = m.spinner.View()
	}
	icon = ui.StatusStyle(t.status).Render(icon)

	dur := ""
	if t.status != ui.StatusRunning && t.status != ui.StatusPending {
		dur = ui.MutedStyle.Render(formatDuration(t.duration))
	}
	name := t.name
	if t.status == ui.StatusFailed || t.status == ui.StatusErrored {
		name = ui.FailureStyle.Render(name)
	}
	parts := []string{icon, name}
	if dur != "" {
		parts = append(parts, dur)
	}
	return strings.Join(parts, "  ")
}

func (m *Model) renderProgress() string {
	if m.totalTests == 0 {
		return ""
	}
	completed := 0
	for _, t := range m.testState {
		if t.status != ui.StatusRunning && t.status != ui.StatusPending {
			completed++
		}
	}
	pct := float64(completed) / float64(m.totalTests)
	if pct > 1 {
		pct = 1
	}
	bar := m.progress.ViewAs(pct)
	wave := ui.MutedStyle.Render(fmt.Sprintf("wave %d/%d · %d/%d tests",
		m.currentWave, m.totalWaves, completed, m.totalTests))
	return lipgloss.JoinHorizontal(lipgloss.Top, bar, "  ", wave)
}

// renderSummary returns one line: "3 passed, 1 failed · 1.85s". No border.
func (m *Model) renderSummary() string {
	if m.summary == nil {
		return ""
	}
	s := m.summary
	parts := []string{ui.SuccessStyle.Render(fmt.Sprintf("%d passed", s.Passed))}
	if s.Failed > 0 {
		parts = append(parts, ui.FailureStyle.Render(fmt.Sprintf("%d failed", s.Failed)))
	}
	if s.Errored > 0 {
		parts = append(parts, ui.FailureStyle.Render(fmt.Sprintf("%d errored", s.Errored)))
	}
	if s.Skipped > 0 {
		parts = append(parts, ui.MutedStyle.Render(fmt.Sprintf("%d skipped", s.Skipped)))
	}
	dur := ui.MutedStyle.Render(formatDuration(s.Duration))
	return strings.Join(parts, ", ") + " · " + dur
}

func (m *Model) renderFailureDetails() string {
	var lines []string
	for _, name := range m.testOrder {
		t := m.testState[name]
		if t.status != ui.StatusFailed && t.status != ui.StatusErrored {
			continue
		}
		lines = append(lines, ui.FailureStyle.Render("✗ "+t.name))
		if t.err != "" {
			lines = append(lines, "  "+ui.MutedStyle.Render(t.err))
		}
		for _, a := range t.failures {
			lines = append(lines,
				"  "+ui.FailureStyle.Render(a.Expression)+"  "+
					ui.MutedStyle.Render(fmt.Sprintf("expected %v, actual %v", a.Expected, a.Actual)),
			)
		}
	}
	return strings.Join(lines, "\n")
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

func failedAssertions(in []engine.AssertionResult) []engine.AssertionResult {
	var out []engine.AssertionResult
	for _, a := range in {
		if !a.Passed {
			out = append(out, a)
		}
	}
	return out
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
