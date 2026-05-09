package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/apiqube/engine"
)

// JSON is an EventHandler that writes each event as a single JSON line.
// Output is NDJSON: one object per line, each carrying a "type" tag and
// the original event payload.
type JSON struct {
	w  io.Writer
	mu sync.Mutex
}

// NewJSON creates a JSON output handler writing to w.
func NewJSON(w io.Writer) *JSON {
	return &JSON{w: w}
}

// Handle marshals the event and writes it as a single line. Errors are
// reported to stderr; the handler never panics.
func (j *JSON) Handle(event engine.Event) {
	j.mu.Lock()
	defer j.mu.Unlock()

	wrapper := struct {
		Type    string      `json:"type"`
		Payload engine.Event `json:"payload"`
	}{
		Type:    event.Type(),
		Payload: event,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Fprintf(os.Stderr, "qube/output: json marshal: %v\n", err)
		return
	}
	if _, err := j.w.Write(append(data, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "qube/output: json write: %v\n", err)
	}
}
