// Package output provides EventHandler implementations for different output formats.
//
// Each format is a separate type that implements engine.EventHandler:
//   - Pretty: colorful terminal output with tables, progress bars, summaries
//   - JSON:   one JSON line per event, suitable for piping to jq or log aggregators
//   - JUnit:  JUnit-compatible XML for CI integration
//   - TAP:    Test Anything Protocol for TAP consumers
package output
