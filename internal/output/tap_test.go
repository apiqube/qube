package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestTAP_PlanAndCases(t *testing.T) {
	var buf bytes.Buffer
	h := NewTAP(&buf)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}

	got := buf.String()

	if !strings.Contains(got, "TAP version 13") {
		t.Error("missing TAP version line")
	}
	if !strings.Contains(got, "1..3") {
		t.Error("missing plan line")
	}
	if !strings.Contains(got, "ok 1 - fetch") {
		t.Error("missing first ok line")
	}
	if !strings.Contains(got, "not ok 2 - create") {
		t.Error("missing second not-ok line")
	}
	if !strings.Contains(got, "not ok 3 - broken") {
		t.Error("missing third not-ok line")
	}
	if !strings.Contains(got, "  ---") {
		t.Error("missing YAML diagnostics block")
	}
	if !strings.Contains(got, "  ...") {
		t.Error("missing YAML diagnostics terminator")
	}
}
