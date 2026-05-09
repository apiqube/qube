package output

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/apiqube/engine"
)

// JUnit is an EventHandler that collects test results and emits a JUnit XML
// report on RunCompleted. Used by CI systems (GitHub Actions, GitLab,
// Jenkins) that natively render JUnit output.
type JUnit struct {
	w       io.Writer
	mu      sync.Mutex
	results []engine.TestResult
	total   engine.RunCompleted
}

// NewJUnit creates a JUnit output handler writing to w.
func NewJUnit(w io.Writer) *JUnit {
	return &JUnit{w: w}
}

// Handle accumulates events and emits the XML on run completion.
func (j *JUnit) Handle(event engine.Event) {
	j.mu.Lock()
	defer j.mu.Unlock()

	switch e := event.(type) {
	case engine.TestCompleted:
		j.results = append(j.results, e.TestResult)
	case engine.RunCompleted:
		j.total = e
		j.emit()
	}
}

func (j *JUnit) emit() {
	suite := xmlTestSuite{
		Name:      "qube",
		Tests:     j.total.Total,
		Failures:  j.total.Failed,
		Errors:    j.total.Errored,
		Skipped:   j.total.Skipped,
		TimeSec:   j.total.Duration.Seconds(),
		TestCases: make([]xmlTestCase, 0, len(j.results)),
	}
	for _, r := range j.results {
		suite.TestCases = append(suite.TestCases, convertCase(r))
	}
	suites := xmlTestSuites{Suites: []xmlTestSuite{suite}}

	if _, err := io.WriteString(j.w, xml.Header); err != nil {
		fmt.Fprintf(os.Stderr, "qube/output: junit header: %v\n", err)
	}
	enc := xml.NewEncoder(j.w)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		fmt.Fprintf(os.Stderr, "qube/output: junit encode: %v\n", err)
	}
	_, _ = io.WriteString(j.w, "\n")
}

func convertCase(r engine.TestResult) xmlTestCase {
	tc := xmlTestCase{
		Name:      r.Name,
		ClassName: r.File,
		TimeSec:   r.Duration.Seconds(),
	}
	switch r.Status {
	case engine.StatusFailed:
		tc.Failure = &xmlMessage{
			Message: failureMessage(r),
			Body:    detailBody(r),
		}
	case engine.StatusErrored:
		tc.Error = &xmlMessage{
			Message: r.Error,
			Body:    detailBody(r),
		}
	case engine.StatusSkipped:
		tc.Skipped = &xmlSkipped{}
	}
	return tc
}

func failureMessage(r engine.TestResult) string {
	for _, a := range r.Assertions {
		if !a.Passed {
			if a.Message != "" {
				return a.Message
			}
			return fmt.Sprintf("%s failed", a.Expression)
		}
	}
	if r.Error != "" {
		return r.Error
	}
	return "test failed"
}

func detailBody(r engine.TestResult) string {
	var b []byte
	for _, a := range r.Assertions {
		if a.Passed {
			continue
		}
		b = fmt.Appendf(b, "%s\n  expected: %v\n  actual:   %v\n  message:  %s\n",
			a.Expression, a.Expected, a.Actual, a.Message)
	}
	if len(b) == 0 && r.Error != "" {
		return r.Error
	}
	return string(b)
}

type xmlTestSuites struct {
	XMLName xml.Name        `xml:"testsuites"`
	Suites  []xmlTestSuite  `xml:"testsuite"`
}

type xmlTestSuite struct {
	XMLName   xml.Name      `xml:"testsuite"`
	Name      string        `xml:"name,attr"`
	Tests     int           `xml:"tests,attr"`
	Failures  int           `xml:"failures,attr"`
	Errors    int           `xml:"errors,attr"`
	Skipped   int           `xml:"skipped,attr"`
	TimeSec   float64       `xml:"time,attr"`
	TestCases []xmlTestCase `xml:"testcase"`
}

type xmlTestCase struct {
	XMLName   xml.Name    `xml:"testcase"`
	Name      string      `xml:"name,attr"`
	ClassName string      `xml:"classname,attr,omitempty"`
	TimeSec   float64     `xml:"time,attr"`
	Failure   *xmlMessage `xml:"failure,omitempty"`
	Error     *xmlMessage `xml:"error,omitempty"`
	Skipped   *xmlSkipped `xml:"skipped,omitempty"`
}

type xmlMessage struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
}

type xmlSkipped struct {
	XMLName xml.Name `xml:"skipped"`
}
