package output

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestJUnit_FullRun(t *testing.T) {
	var buf bytes.Buffer
	h := NewJUnit(&buf)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}

	out := buf.String()
	if !strings.HasPrefix(out, xml.Header) {
		t.Error("output should start with the XML header")
	}

	var suites xmlTestSuites
	if err := xml.Unmarshal([]byte(strings.TrimPrefix(out, xml.Header)), &suites); err != nil {
		t.Fatalf("output is not parseable JUnit: %v\n%s", err, out)
	}
	if len(suites.Suites) != 1 {
		t.Fatalf("got %d suites; want 1", len(suites.Suites))
	}
	suite := suites.Suites[0]
	if suite.Tests != 3 {
		t.Errorf("tests = %d; want 3", suite.Tests)
	}
	if suite.Failures != 1 {
		t.Errorf("failures = %d; want 1", suite.Failures)
	}
	if suite.Errors != 1 {
		t.Errorf("errors = %d; want 1", suite.Errors)
	}

	if len(suite.TestCases) != 3 {
		t.Fatalf("got %d testcases; want 3", len(suite.TestCases))
	}
	pass, fail, errd := suite.TestCases[0], suite.TestCases[1], suite.TestCases[2]
	if pass.Failure != nil || pass.Error != nil {
		t.Errorf("passing test should have no failure/error: %+v", pass)
	}
	if fail.Failure == nil {
		t.Error("failed test should have <failure>")
	}
	if fail.Failure != nil && !strings.Contains(fail.Failure.Body, "expected") {
		t.Errorf("failure body missing expected/actual: %q", fail.Failure.Body)
	}
	if errd.Error == nil {
		t.Error("errored test should have <error>")
	}
	if errd.Error != nil && errd.Error.Message != "connection refused" {
		t.Errorf("error message = %q; want 'connection refused'", errd.Error.Message)
	}
}

func TestJUnit_EmptyRun(t *testing.T) {
	var buf bytes.Buffer
	h := NewJUnit(&buf)
	// Just emit the run-completion event with no tests in between.
	h.Handle(canonicalEvents()[0])
	for _, e := range canonicalEvents() {
		if _, ok := e.(interface{ Type() string }); ok {
			if e.Type() == "RunCompleted" {
				h.Handle(e)
			}
		}
	}
	if buf.Len() == 0 {
		t.Error("empty run should still emit XML")
	}
}
