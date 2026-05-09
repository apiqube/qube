package run

import (
	"strings"
	"testing"
	"time"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/ui"
)

func feed(m *Model, events ...engine.Event) {
	for _, e := range events {
		m.handleEngineEvent(e)
	}
}

func TestModel_RunStartedPopulatesHeader(t *testing.T) {
	m := New()
	feed(m, engine.RunStarted{Files: []string{"a.yaml"}, TotalTests: 5, TotalWaves: 2})

	if m.totalTests != 5 {
		t.Errorf("totalTests = %d; want 5", m.totalTests)
	}
	if m.totalWaves != 2 {
		t.Errorf("totalWaves = %d; want 2", m.totalWaves)
	}
	if len(m.files) != 1 {
		t.Errorf("files len = %d; want 1", len(m.files))
	}
}

func TestModel_TestStartedAddsToOrderAndStateRunning(t *testing.T) {
	m := New()
	feed(m, engine.TestStarted{Name: "fetch"})
	if len(m.testOrder) != 1 || m.testOrder[0] != "fetch" {
		t.Errorf("order wrong: %v", m.testOrder)
	}
	if m.testState["fetch"].status != ui.StatusRunning {
		t.Errorf("status should be Running, got %v", m.testState["fetch"].status)
	}
}

func TestModel_TestCompletedUpdatesStatusAndDuration(t *testing.T) {
	m := New()
	feed(m,
		engine.TestStarted{Name: "fetch"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:     "fetch",
			Status:   engine.StatusPassed,
			Duration: 25 * time.Millisecond,
		}},
	)
	if m.testState["fetch"].status != ui.StatusPassed {
		t.Errorf("status not updated: %v", m.testState["fetch"].status)
	}
	if m.testState["fetch"].duration != 25*time.Millisecond {
		t.Errorf("duration not stored: %v", m.testState["fetch"].duration)
	}
}

func TestModel_TestCompletedCapturesFailureDetails(t *testing.T) {
	m := New()
	feed(m,
		engine.TestStarted{Name: "create"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:   "create",
			Status: engine.StatusFailed,
			Assertions: []engine.AssertionResult{
				{Expression: "status", Passed: false, Expected: 201, Actual: 500},
				{Expression: "body.id", Passed: true},
			},
		}},
	)
	entry := m.testState["create"]
	if entry.status != ui.StatusFailed {
		t.Errorf("status wrong: %v", entry.status)
	}
	if len(entry.failures) != 1 {
		t.Errorf("expected 1 failure, got %d", len(entry.failures))
	}
}

func TestModel_RunCompletedStoresSummary(t *testing.T) {
	m := New()
	feed(m, engine.RunCompleted{Total: 3, Passed: 2, Failed: 1, Duration: 100 * time.Millisecond})
	if m.summary == nil {
		t.Fatal("summary should be set")
	}
	if m.summary.Passed != 2 {
		t.Errorf("summary.Passed = %d; want 2", m.summary.Passed)
	}
}

func TestModel_ViewLiveContainsKeyFragments(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	m := New()
	feed(m,
		engine.RunStarted{Files: []string{"a.yaml"}, TotalTests: 2, TotalWaves: 1},
		engine.TestStarted{Name: "first"},
	)
	v := m.View()
	if !strings.Contains(v, "qube run") {
		t.Errorf("view missing brand: %q", v)
	}
	if !strings.Contains(v, "first") {
		t.Errorf("view missing test name: %q", v)
	}
}

func TestModel_ViewSummaryRendersWhenFinished(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	m := New()
	feed(m,
		engine.RunStarted{Files: []string{"a.yaml"}, TotalTests: 1, TotalWaves: 1},
		engine.TestStarted{Name: "fail"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:   "fail",
			Status: engine.StatusFailed,
			Assertions: []engine.AssertionResult{
				{Expression: "status", Passed: false, Expected: 200, Actual: 500},
			},
		}},
		engine.RunCompleted{Total: 1, Failed: 1, Duration: time.Millisecond},
	)
	m.finished = true

	v := m.View()
	if !strings.Contains(v, "1 failed") {
		t.Errorf("summary missing failed count: %q", v)
	}
	if !strings.Contains(v, "fail") {
		t.Errorf("summary missing failure detail: %q", v)
	}
}

func TestMapStatus(t *testing.T) {
	cases := map[engine.TestStatus]ui.Status{
		engine.StatusPassed:  ui.StatusPassed,
		engine.StatusFailed:  ui.StatusFailed,
		engine.StatusSkipped: ui.StatusSkipped,
		engine.StatusErrored: ui.StatusErrored,
	}
	for in, want := range cases {
		if got := mapStatus(in); got != want {
			t.Errorf("mapStatus(%v) = %v; want %v", in, got, want)
		}
	}
}

func TestFormatDurationLocal(t *testing.T) {
	cases := map[time.Duration]string{
		1 * time.Microsecond: "1µs",
		2 * time.Millisecond: "2ms",
		3 * time.Second:      "3.00s",
	}
	for d, want := range cases {
		if got := formatDuration(d); got != want {
			t.Errorf("formatDuration(%v) = %q; want %q", d, got, want)
		}
	}
}
