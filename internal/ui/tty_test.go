package ui

import (
	"bytes"
	"testing"
)

func TestIsTTY_NonFile(t *testing.T) {
	if IsTTY(&bytes.Buffer{}) {
		t.Error("buffer should not be a TTY")
	}
}

func TestIsInteractive_HonorsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if IsInteractive() {
		t.Error("NO_COLOR should disable interactive mode")
	}
}

func TestIsInteractive_HonorsCI(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("CI", "true")
	if IsInteractive() {
		t.Error("CI=true should disable interactive mode")
	}
}

func TestIsInteractive_HonorsDumbTerm(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("CI", "")
	t.Setenv("TERM", "dumb")
	if IsInteractive() {
		t.Error("TERM=dumb should disable interactive mode")
	}
}
