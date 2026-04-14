package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "qube",
	Short: "Declarative API testing CLI",
	Long: `qube — test HTTP, gRPC, GraphQL, WebSocket and more with one YAML format.

Write tests in YAML, run them with auto-parallelism, and get structured results.
Features automatic dependency detection, cross-test data flow, and plugin support.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns any error that occurred.
func Execute() error {
	return rootCmd.Execute()
}

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
