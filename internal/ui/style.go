// Package ui holds the qube CLI's shared style palette, status icons,
// log writer, and TTY detection. Output is log-line oriented
// (`timestamp LEVEL message`) — no rounded borders, no boxed cards.
package ui

import "github.com/charmbracelet/lipgloss"

// Adaptive colors. Each renders correctly on dark and light terminals.
var (
	ColorBrand   = lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#5BC8E8"}
	ColorSuccess = lipgloss.AdaptiveColor{Light: "#0E7C3A", Dark: "#5DD68A"}
	ColorWarn    = lipgloss.AdaptiveColor{Light: "#B7791F", Dark: "#F6C667"}
	ColorFailure = lipgloss.AdaptiveColor{Light: "#C0392B", Dark: "#FF6B6B"}
	ColorMuted   = lipgloss.AdaptiveColor{Light: "#6C6C6C", Dark: "#A0A0A0"}
	ColorAccent  = lipgloss.AdaptiveColor{Light: "#7C2D8B", Dark: "#C084FC"}
	ColorDebug   = lipgloss.AdaptiveColor{Light: "#3060A0", Dark: "#7AA8E8"}
)

// Text-only styles (foreground colors). No borders, no backgrounds — keeps
// the output looking like a structured log instead of a UI dashboard.
var (
	BrandStyle   = lipgloss.NewStyle().Foreground(ColorBrand).Bold(true)
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
	FailureStyle = lipgloss.NewStyle().Foreground(ColorFailure).Bold(true)
	WarnStyle    = lipgloss.NewStyle().Foreground(ColorWarn)
	MutedStyle   = lipgloss.NewStyle().Foreground(ColorMuted)
	AccentStyle  = lipgloss.NewStyle().Foreground(ColorAccent)
	DebugStyle   = lipgloss.NewStyle().Foreground(ColorDebug)

	// timestampStyle dims the leading clock value in log lines.
	timestampStyle = MutedStyle

	// TableHeaderStyle is the bold, brand-colored header row of any table
	// the CLI prints (plugin list, etc).
	TableHeaderStyle = lipgloss.NewStyle().Foreground(ColorBrand).Bold(true)
)

// levelStyle returns the foreground style for one log Level.
func levelStyle(l Level) lipgloss.Style {
	switch l {
	case LevelDebug:
		return DebugStyle
	case LevelInfo:
		return BrandStyle
	case LevelSuccess:
		return SuccessStyle
	case LevelWarn:
		return WarnStyle
	case LevelError:
		return FailureStyle
	}
	return BrandStyle
}

// StatusStyle returns the style suited for a Status icon.
func StatusStyle(s Status) lipgloss.Style {
	switch s {
	case StatusPassed:
		return SuccessStyle
	case StatusFailed, StatusErrored:
		return FailureStyle
	case StatusSkipped:
		return MutedStyle
	case StatusRunning:
		return AccentStyle
	}
	return MutedStyle
}
