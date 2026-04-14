package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkFlags struct {
	configPath string
	pluginDir  string
}

var checkCmd = &cobra.Command{
	Use:   "check [path...]",
	Short: "Validate test manifests without executing them",
	Long: `Check validates YAML syntax, known/unknown fields, plugin requirements,
unresolved template references, and dependency cycles.

Exits with non-zero status if any errors are found.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// 1. Discover config and plugins (same as run)
		// 2. Build engine with config
		// 3. Call eng.Check(ctx, engine.FromPaths(paths...))
		// 4. Render validation errors (group by file, show line numbers)
		// 5. Exit non-zero if any errors (warnings are OK)
		return fmt.Errorf("not implemented")
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkFlags.configPath, "config", "c", "", "path to .qube.yaml")
	checkCmd.Flags().StringVar(&checkFlags.pluginDir, "plugins", "", "plugin directory")
}
