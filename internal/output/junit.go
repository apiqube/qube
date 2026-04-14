package output

import (
	"io"

	"github.com/apiqube/engine"
)

// JUnit is an EventHandler that collects results and writes JUnit XML
// at the end of the run. Used for CI integration (GitLab, GitHub, Jenkins).
type JUnit struct {
	w       io.Writer
	results []engine.TestResult
}

// NewJUnit creates a new JUnit output handler writing to w.
func NewJUnit(w io.Writer) *JUnit {
	return &JUnit{w: w}
}

// Handle collects TestCompleted events and writes the final XML on RunCompleted.
func (j *JUnit) Handle(event engine.Event) {
	// TODO: implementation
	//
	// case engine.TestCompleted:
	//   j.results = append(j.results, e.TestResult)
	//
	// case engine.RunCompleted:
	//   Build JUnit XML structure:
	//     <testsuites>
	//       <testsuite name="qube" tests="N" failures="F" errors="E" time="T">
	//         <testcase name="..." classname="..." time="...">
	//           <failure message="..."><![CDATA[...]]></failure>
	//         </testcase>
	//         ...
	//       </testsuite>
	//     </testsuites>
	//   Write to j.w
}
