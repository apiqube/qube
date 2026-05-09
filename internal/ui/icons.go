package ui

import "os"

// Status enumerates the per-test states the CLI renders.
type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusPassed
	StatusFailed
	StatusSkipped
	StatusErrored
)

// String returns the canonical name (matches engine.TestStatus output).
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusPassed:
		return "passed"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	case StatusErrored:
		return "errored"
	}
	return "unknown"
}

// Icon returns the styled glyph for a status. Falls back to ASCII when the
// terminal is plain (NO_COLOR set, TERM=dumb).
func (s Status) Icon() string {
	if asciiOnly() {
		return s.asciiIcon()
	}
	return s.unicodeIcon()
}

func (s Status) unicodeIcon() string {
	switch s {
	case StatusPending:
		return "○"
	case StatusRunning:
		return "▶"
	case StatusPassed:
		return "✓"
	case StatusFailed:
		return "✗"
	case StatusSkipped:
		return "⏭"
	case StatusErrored:
		return "⚠"
	}
	return "?"
}

func (s Status) asciiIcon() string {
	switch s {
	case StatusPending:
		return "."
	case StatusRunning:
		return ">"
	case StatusPassed:
		return "PASS"
	case StatusFailed:
		return "FAIL"
	case StatusSkipped:
		return "SKIP"
	case StatusErrored:
		return "ERR "
	}
	return "?"
}

// asciiOnly reports whether icons should fall back to ASCII.
func asciiOnly() bool {
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	if t := os.Getenv("TERM"); t == "dumb" || t == "" {
		// dumb terminals lack styling support; "" is common in cron / non-interactive shells.
		return t == "dumb"
	}
	return false
}
