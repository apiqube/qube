// Package commands defines the qube CLI command tree, built on cobra and
// dressed by fang. Each command lives in its own file; this package wires
// them together and exposes Root() for cmd/main to drive.
package commands

import (
	"github.com/spf13/cobra"

	"github.com/apiqube/qube/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "qube",
	Short: "Declarative API testing CLI",
	Long: `qube — test HTTP, gRPC, GraphQL, WebSocket and more with one YAML format.

Write tests in YAML, run them with auto-parallelism, and get structured results.
Features automatic dependency detection, cross-test data flow, and plugin support.`,
	Version:       version.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Root returns the cobra root command. fang.Execute consumes it from
// cmd/main.go.
func Root() *cobra.Command { return rootCmd }

func init() {
	rootCmd.AddCommand(
		runCmd,
		checkCmd,
		initCmd,
		generateCmd,
		pluginCmd,
		versionCmd,
	)
}
