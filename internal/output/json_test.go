package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSON_NDJSONShape(t *testing.T) {
	var buf bytes.Buffer
	h := NewJSON(&buf)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != len(canonicalEvents()) {
		t.Fatalf("expected %d lines, got %d", len(canonicalEvents()), len(lines))
	}
	for i, line := range lines {
		var wrapper struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal([]byte(line), &wrapper); err != nil {
			t.Fatalf("line %d: invalid JSON: %v\n%s", i, err, line)
		}
		if wrapper.Type == "" {
			t.Errorf("line %d: missing type field", i)
		}
	}
}

func TestJSON_TypeOrder(t *testing.T) {
	var buf bytes.Buffer
	h := NewJSON(&buf)
	for _, e := range canonicalEvents() {
		h.Handle(e)
	}

	want := []string{
		"RunStarted",
		"TestStarted", "TestCompleted",
		"TestStarted", "TestCompleted",
		"TestStarted", "TestCompleted",
		"RunCompleted",
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	for i, line := range lines {
		var wrapper struct {
			Type string `json:"type"`
		}
		_ = json.Unmarshal([]byte(line), &wrapper)
		if wrapper.Type != want[i] {
			t.Errorf("line %d: type = %q; want %q", i, wrapper.Type, want[i])
		}
	}
}
