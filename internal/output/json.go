package output

import (
	"io"

	"github.com/apiqube/engine"
)

// JSON is an EventHandler that writes each event as a single JSON line.
// Useful for piping to jq, log aggregators, or external processors.
type JSON struct {
	w io.Writer
}

// NewJSON creates a new JSON output handler writing to w.
func NewJSON(w io.Writer) *JSON {
	return &JSON{w: w}
}

// Handle marshals the event to JSON and writes it as one line (NDJSON).
func (j *JSON) Handle(event engine.Event) {
	// TODO: implementation
	//
	// 1. Build wrapper with "type": event.Type() and "payload": event
	// 2. Marshal to JSON
	// 3. Write line with trailing newline
	// 4. On error, write to stderr (never panic on output)
}
