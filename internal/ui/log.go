package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Level classifies the kind of log line.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelSuccess
	LevelWarn
	LevelError
)

// Width of the level column when padded for alignment.
const levelWidth = 7

// String returns the canonical name shown in a log line.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelSuccess:
		return "SUCCESS"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "INFO"
}

// Logger writes structured log lines: `timestamp LEVEL message`.
//
// Output looks like
//
//	14:30:45 INFO    qube run · 4 tests · 1 file · 1 wave
//	14:30:46 SUCCESS ✓ POST /post (124ms)
//
// Continuation lines (after `\n` in the message) are indented to align with
// the message column. The default Logger writes to stderr; tests substitute
// their own writer via With.
type Logger struct {
	w        io.Writer
	mu       sync.Mutex
	clock    func() time.Time
	verbose  bool
	noColor  bool
}

// DefaultLogger is the package-level Logger used by the convenience helpers
// (Info / Success / Warn / Error / Debug). Commands typically use the
// helpers directly; tests can swap the underlying writer with WithWriter.
var DefaultLogger = &Logger{
	w:       os.Stderr,
	clock:   time.Now,
	noColor: asciiOnly(),
}

// WithWriter returns a copy of l writing to w. Used by tests.
func (l *Logger) WithWriter(w io.Writer) *Logger {
	return &Logger{w: w, clock: l.clock, verbose: l.verbose, noColor: l.noColor}
}

// SetVerbose toggles emission of debug-level lines.
func (l *Logger) SetVerbose(v bool) { l.verbose = v }

// SetWriter replaces the destination writer in place. Used by tests; commands
// should leave the default (os.Stderr).
func (l *Logger) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.w = w
}

// IsVerbose reports whether debug lines are emitted.
func (l *Logger) IsVerbose() bool { return l.verbose }

// Log writes one log entry. msg may contain newlines; subsequent lines are
// indented to align with the message column.
func (l *Logger) Log(level Level, msg string) {
	if level == LevelDebug && !l.verbose {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := timestampStyle.Render(l.clock().Format("15:04:05"))
	levelText := fmt.Sprintf("%-*s", levelWidth, level.String())
	levelStyled := levelStyle(level).Render(levelText)

	indent := len("15:04:05") + 1 + levelWidth + 1
	body := indentMultiline(msg, indent)

	fmt.Fprintf(l.w, "%s %s %s\n", timestamp, levelStyled, body)
}

// Logf is the printf-style sibling of Log.
func (l *Logger) Logf(level Level, format string, args ...any) {
	l.Log(level, fmt.Sprintf(format, args...))
}

// indentMultiline prefixes every line after the first with `width` spaces so
// continuation text aligns with the message column.
func indentMultiline(msg string, width int) string {
	if !strings.Contains(msg, "\n") {
		return msg
	}
	lines := strings.Split(msg, "\n")
	pad := strings.Repeat(" ", width)
	for i := 1; i < len(lines); i++ {
		lines[i] = pad + lines[i]
	}
	return strings.Join(lines, "\n")
}

// Convenience helpers below — they delegate to DefaultLogger.

// Debug writes a DEBUG line; suppressed unless --verbose was set.
func Debug(msg string) { DefaultLogger.Log(LevelDebug, msg) }

// Debugf is the printf-style sibling of Debug.
func Debugf(format string, args ...any) { DefaultLogger.Logf(LevelDebug, format, args...) }

// Info writes an INFO line.
func Info(msg string) { DefaultLogger.Log(LevelInfo, msg) }

// Infof is the printf-style sibling of Info.
func Infof(format string, args ...any) { DefaultLogger.Logf(LevelInfo, format, args...) }

// Success writes a SUCCESS line (green).
func Success(msg string) { DefaultLogger.Log(LevelSuccess, msg) }

// Successf is the printf-style sibling of Success.
func Successf(format string, args ...any) { DefaultLogger.Logf(LevelSuccess, format, args...) }

// Warn writes a WARN line (yellow).
func Warn(msg string) { DefaultLogger.Log(LevelWarn, msg) }

// Warnf is the printf-style sibling of Warn.
func Warnf(format string, args ...any) { DefaultLogger.Logf(LevelWarn, format, args...) }

// Err writes an ERROR line (red, bold). Named Err to avoid colliding with the
// `error` builtin in user code that imports this package as `ui`.
func Err(msg string) { DefaultLogger.Log(LevelError, msg) }

// Errf is the printf-style sibling of Err.
func Errf(format string, args ...any) { DefaultLogger.Logf(LevelError, format, args...) }
