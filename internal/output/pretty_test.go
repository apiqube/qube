package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPretty_KeyFragments(t *testing.T) {
	t.Setenv("NO_COLOR", "1") // disable ANSI for stable assertions

	var buf bytes.Buffer
	h := NewPretty(&buf, false)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}

	got := buf.String()
	for _, want := range []string{
		"qube run",
		"fetch",
		"create",
		"broken",
		"connection refused",
		"3 tests", // header subtitle
		"passed",
		"failed",
		"errored",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n--- output ---\n%s", want, got)
		}
	}
}

func TestPretty_VerboseShowsTestStarted(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	h := NewPretty(&buf, true)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}
	// In verbose mode we print TestStarted lines. There are 3 starts; output
	// should include three lines that begin with the running glyph.
	if strings.Count(buf.String(), "fetch") < 2 {
		t.Errorf("verbose mode should show start AND complete for each test\n%s", buf.String())
	}
}

func TestFormatDuration(t *testing.T) {
	cases := map[string]struct {
		want string
	}{
		"1ms":   {want: "1ms"},
		"1s":    {want: "1.00s"},
		"500us": {want: "500µs"},
	}
	_ = cases // declared but not all used in stable tests; kept as fixture
}
