package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qube",
	Short: "Declarative API testing CLI",
	Long:  "qube — test HTTP, gRPC, GraphQL, WebSocket and more with one YAML format.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	version = "dev"
	commit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "qube %s (%s)\n", version, commit)
	},
}
