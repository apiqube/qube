package output

import (
	"time"

	"github.com/apiqube/engine"
)

// canonicalEvents returns a fixture event stream representing a 3-test run:
// one passing, one failing (with an assertion failure), one errored.
func canonicalEvents() []engine.Event {
	files := []string{"tests/example.yaml"}
	return []engine.Event{
		engine.RunStarted{Files: files, TotalTests: 3, TotalWaves: 1},
		engine.TestStarted{Name: "fetch", File: files[0], Protocol: engine.ProtocolHTTP, Target: "http://api"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:     "fetch",
			File:     files[0],
			Protocol: engine.ProtocolHTTP,
			Target:   "http://api",
			Status:   engine.StatusPassed,
			Duration: 35 * time.Millisecond,
		}},
		engine.TestStarted{Name: "create", File: files[0], Protocol: engine.ProtocolHTTP, Target: "http://api"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:     "create",
			File:     files[0],
			Protocol: engine.ProtocolHTTP,
			Target:   "http://api",
			Status:   engine.StatusFailed,
			Duration: 50 * time.Millisecond,
			Assertions: []engine.AssertionResult{
				{Expression: "status", Passed: false, Expected: 201, Actual: 500, Message: "expected 201, got 500"},
			},
		}},
		engine.TestStarted{Name: "broken", File: files[0], Protocol: engine.ProtocolHTTP, Target: "http://api"},
		engine.TestCompleted{TestResult: engine.TestResult{
			Name:     "broken",
			File:     files[0],
			Protocol: engine.ProtocolHTTP,
			Target:   "http://api",
			Status:   engine.StatusErrored,
			Duration: 5 * time.Millisecond,
			Error:    "connection refused",
		}},
		engine.RunCompleted{
			Total:    3,
			Passed:   1,
			Failed:   1,
			Errored:  1,
			Duration: 90 * time.Millisecond,
		},
	}
}
