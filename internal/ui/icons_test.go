package ui

import (
	"testing"
)

func TestStatusString(t *testing.T) {
	cases := map[Status]string{
		StatusPending: "pending",
		StatusRunning: "running",
		StatusPassed:  "passed",
		StatusFailed:  "failed",
		StatusSkipped: "skipped",
		StatusErrored: "errored",
		Status(99):    "unknown",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q; want %q", s, got, want)
		}
	}
}

func TestStatusIcon_Unicode(t *testing.T) {
	// Ensure NO_COLOR is unset for this test; explicitly set TERM to a real one.
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")

	cases := []struct {
		s    Status
		want string
	}{
		{StatusPending, "○"},
		{StatusRunning, "▶"},
		{StatusPassed, "✓"},
		{StatusFailed, "✗"},
		{StatusSkipped, "⏭"},
		{StatusErrored, "⚠"},
	}
	for _, c := range cases {
		if got := c.s.Icon(); got != c.want {
			t.Errorf("%v.Icon() = %q; want %q", c.s, got, c.want)
		}
	}
}

func TestStatusIcon_ASCIIFallback_NoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	if got := StatusPassed.Icon(); got != "PASS" {
		t.Errorf("under NO_COLOR, passed icon = %q; want PASS", got)
	}
	if got := StatusFailed.Icon(); got != "FAIL" {
		t.Errorf("under NO_COLOR, failed icon = %q; want FAIL", got)
	}
}

func TestStatusIcon_ASCIIFallback_DumbTerm(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")

	if got := StatusPassed.Icon(); got != "PASS" {
		t.Errorf("under TERM=dumb, passed icon = %q; want PASS", got)
	}
}

func TestStatusIcon_Unknown(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	if got := Status(999).Icon(); got != "?" {
		t.Errorf("unknown status icon = %q; want '?'", got)
	}
}
