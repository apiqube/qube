package ui

import (
	"io"
	"os"

	"golang.org/x/term"
)

// IsTTY reports whether the given writer points at a terminal capable of
// styling. Anything not backed by an *os.File is treated as non-TTY.
func IsTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

// IsInteractive reports whether the CLI should drive a full TUI.
//
// Returns false when stdout is not a terminal, when NO_COLOR is set, or when
// CI environment variables are present. Used to pick between Bubble Tea and
// progressive lipgloss output.
func IsInteractive() bool {
	if !IsTTY(os.Stdout) {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("CI") != "" {
		return false
	}
	if t := os.Getenv("TERM"); t == "dumb" {
		return false
	}
	return true
}
