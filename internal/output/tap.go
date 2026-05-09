package output

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/apiqube/engine"
)

// TAP is an EventHandler that emits Test Anything Protocol (version 13) output.
// One `1..N` plan, one `ok N - name` line per test, with a YAML diagnostic
// block under failed/errored cases.
type TAP struct {
	w       io.Writer
	mu      sync.Mutex
	planned bool
	idx     int
}

// NewTAP creates a TAP output handler writing to w.
func NewTAP(w io.Writer) *TAP {
	return &TAP{w: w}
}

// Handle emits TAP lines as events arrive.
func (t *TAP) Handle(event engine.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch e := event.(type) {
	case engine.RunStarted:
		fmt.Fprintln(t.w, "TAP version 13")
		fmt.Fprintf(t.w, "1..%d\n", e.TotalTests)
		t.planned = true
	case engine.TestCompleted:
		t.idx++
		t.writeTestCase(e.TestResult)
	}
}

func (t *TAP) writeTestCase(r engine.TestResult) {
	prefix := "ok"
	suffix := ""
	switch r.Status {
	case engine.StatusFailed, engine.StatusErrored:
		prefix = "not ok"
	case engine.StatusSkipped:
		suffix = " # SKIP"
	}
	if _, err := fmt.Fprintf(t.w, "%s %d - %s%s\n", prefix, t.idx, r.Name, suffix); err != nil {
		fmt.Fprintf(os.Stderr, "qube/output: tap write: %v\n", err)
	}

	if r.Status == engine.StatusFailed || r.Status == engine.StatusErrored {
		t.writeYAMLBlock(r)
	}
}

func (t *TAP) writeYAMLBlock(r engine.TestResult) {
	fmt.Fprintln(t.w, "  ---")
	fmt.Fprintf(t.w, "  status: %s\n", r.Status)
	fmt.Fprintf(t.w, "  duration_ms: %d\n", r.Duration.Milliseconds())
	if r.Protocol != "" {
		fmt.Fprintf(t.w, "  protocol: %s\n", r.Protocol)
	}
	if r.Target != "" {
		fmt.Fprintf(t.w, "  target: %s\n", r.Target)
	}
	if r.Error != "" {
		fmt.Fprintf(t.w, "  error: %q\n", r.Error)
	}
	failed := failedAssertions(r.Assertions)
	if len(failed) > 0 {
		fmt.Fprintln(t.w, "  failures:")
		for _, a := range failed {
			fmt.Fprintf(t.w, "    - expression: %q\n", a.Expression)
			fmt.Fprintf(t.w, "      expected:   %v\n", a.Expected)
			fmt.Fprintf(t.w, "      actual:     %v\n", a.Actual)
			if a.Message != "" {
				fmt.Fprintf(t.w, "      message:    %q\n", a.Message)
			}
		}
	}
	fmt.Fprintln(t.w, "  ...")
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
