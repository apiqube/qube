package commands

import (
	"github.com/spf13/cobra"
)

var generateFlags struct {
	from      string
	output    string
	protoFile string
	format    string
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate tests from API specifications (not yet implemented)",
	Long: `Generate creates qube test files from existing API specifications.

Planned sources:
  OpenAPI/Swagger (JSON or YAML)  --from <url|path>
  Protobuf                         --proto <file>
  HAR (HTTP Archive)               --from <file.har>
  Postman Collection               --from <collection.json>

Generated tests will be written to the directory specified by --output.

This command is on the v1.x roadmap; v1.0 ships with run, check, init,
version, and plugin list. Run "qube --help" to see what's available now.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		printRoadmapStub(cmd.OutOrStdout(), "qube generate",
			"OpenAPI/Protobuf/HAR/Postman parsers and CRUD-chain inference are scoped for a follow-up plan.")
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&generateFlags.from, "from", "", "source: URL or file path")
	generateCmd.Flags().StringVarP(&generateFlags.output, "output", "o", "tests/", "output directory")
	generateCmd.Flags().StringVar(&generateFlags.protoFile, "proto", "", "protobuf file for gRPC generation")
	generateCmd.Flags().StringVar(&generateFlags.format, "format", "compact", "output syntax level: full, compact, oneliner")
}
