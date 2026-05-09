// Package ui holds the qube CLI's shared style palette, status icons, and
// TTY detection. Every command renders through this package so the visual
// language stays consistent across subcommands.
package ui

import "github.com/charmbracelet/lipgloss"

// Palette of named colors used across the CLI. Each entry uses an
// AdaptiveColor so the same code looks right on dark and light terminals.
var (
	ColorBrand   = lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#5BC8E8"}
	ColorSuccess = lipgloss.AdaptiveColor{Light: "#0E7C3A", Dark: "#5DD68A"}
	ColorFailure = lipgloss.AdaptiveColor{Light: "#C0392B", Dark: "#FF6B6B"}
	ColorWarn    = lipgloss.AdaptiveColor{Light: "#B7791F", Dark: "#F6C667"}
	ColorMuted   = lipgloss.AdaptiveColor{Light: "#6C6C6C", Dark: "#A0A0A0"}
	ColorAccent  = lipgloss.AdaptiveColor{Light: "#7C2D8B", Dark: "#C084FC"}
	ColorBorder  = lipgloss.AdaptiveColor{Light: "#D0D0D0", Dark: "#3A3A3A"}
)

// Reusable component styles. Commands compose these — they should never
// instantiate `lipgloss.NewStyle()` directly.
var (
	// Brand is for the header banner / qube branding text.
	Brand = lipgloss.NewStyle().Foreground(ColorBrand).Bold(true)

	// Success / Failure / Warn / Muted are short-form text styles.
	Success = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	Failure = lipgloss.NewStyle().Foreground(ColorFailure).Bold(true)
	Warn    = lipgloss.NewStyle().Foreground(ColorWarn).Bold(true)
	Muted   = lipgloss.NewStyle().Foreground(ColorMuted)
	Accent  = lipgloss.NewStyle().Foreground(ColorAccent)

	// Header is the top-of-output title block (e.g. "qube run").
	Header = lipgloss.NewStyle().
		Foreground(ColorBrand).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(ColorBorder).
		Padding(0, 1).
		MarginBottom(1)

	// Badge wraps a short label like a pill (status, mode, count).
	Badge = lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true)

	// Card frames a group of content with a soft border and padding.
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 2)

	// SummaryCard is the run-end totals block; leans on Card with extra accent.
	SummaryCard = Card.
			BorderForeground(ColorBrand).
			Padding(1, 3)

	// ErrorBlock frames a validation error or panic message.
	ErrorBlock = Card.
			BorderForeground(ColorFailure)

	// TableHeader is the bold header row of any table we render.
	TableHeader = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Bold(true).
			Padding(0, 1)

	// TableCell is a default cell style.
	TableCell = lipgloss.NewStyle().Padding(0, 1)

	// TableCellFaint dims a value (e.g. file paths next to test names).
	TableCellFaint = TableCell.Foreground(ColorMuted)
)

// StatusStyle returns a style suited for the given status icon.
func StatusStyle(s Status) lipgloss.Style {
	switch s {
	case StatusPassed:
		return Success
	case StatusFailed, StatusErrored:
		return Failure
	case StatusSkipped:
		return Muted
	case StatusRunning:
		return Accent
	}
	return Muted
}
