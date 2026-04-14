package commands

import (
	"fmt"

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
	Short: "Generate tests from API specifications",
	Long: `Generate creates qube test files from existing API specifications.

Supported sources:
  OpenAPI/Swagger (JSON or YAML)  --from <url|path>
  Protobuf                         --proto <file>
  HAR (HTTP Archive)               --from <file.har>
  Postman Collection               --from <collection.json>

Generated tests are written to the directory specified by --output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// 1. Determine source type from --from URL/extension or --proto
		// 2. Look up generator plugin for the source type
		// 3. Load source via plugin
		// 4. Generate tests with inferred CRUD chains:
		//    - POST + GET/{id} + PUT/{id} + DELETE/{id} → CRUD scenario
		//    - Response schemas → assertion hints
		//    - Required fields → fake.* templates
		// 5. Write to output directory
		return fmt.Errorf("not implemented")
	},
}

func init() {
	generateCmd.Flags().StringVar(&generateFlags.from, "from", "", "source: URL or file path")
	generateCmd.Flags().StringVarP(&generateFlags.output, "output", "o", "tests/", "output directory")
	generateCmd.Flags().StringVar(&generateFlags.protoFile, "proto", "", "protobuf file for gRPC generation")
	generateCmd.Flags().StringVar(&generateFlags.format, "format", "compact", "output syntax level: full, compact, oneliner")
}
