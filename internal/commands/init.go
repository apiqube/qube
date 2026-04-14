package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initFlags struct {
	interactive bool
	swaggerURL  string
	targetURL   string
	force       bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new qube project",
	Long: `Init creates starter files for a new qube project:
  .qube.yaml         project configuration
  tests/example.yaml simple example test

Use --interactive to answer questions and optionally generate tests
from an existing OpenAPI/Swagger specification.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// Silent mode (default):
		//   1. Create .qube.yaml with default targets and empty plugins list
		//   2. Create tests/example.yaml with a minimal smoke test
		//   3. Print "next steps" message
		//
		// Interactive mode (--interactive):
		//   1. Launch tui.RunInitWizard()
		//      - Ask target URL
		//      - Ask if there's an OpenAPI spec
		//      - If yes, ask for spec URL/path
		//      - Ask which plugins to enable
		//   2. If swagger provided: invoke generate plugin to create tests
		//   3. Otherwise: create example test
		//
		// --force overwrites existing files.
		return fmt.Errorf("not implemented")
	},
}

func init() {
	initCmd.Flags().BoolVarP(&initFlags.interactive, "interactive", "i", false, "ask questions interactively")
	initCmd.Flags().StringVar(&initFlags.swaggerURL, "from-swagger", "", "generate tests from OpenAPI/Swagger spec")
	initCmd.Flags().StringVar(&initFlags.targetURL, "target", "", "default target URL")
	initCmd.Flags().BoolVar(&initFlags.force, "force", false, "overwrite existing files")
}
