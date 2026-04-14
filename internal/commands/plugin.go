package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage qube plugins",
	Long: `Plugin manages WASM plugins that extend qube with new protocols, reporters,
generators, and hooks.`,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name|url|path>",
	Short: "Install a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// 1. Determine source type:
		//    - Bare name → fetch from official registry
		//    - URL → download .wasm file
		//    - Local path → copy file
		// 2. Verify it's a valid WASM module with plugin_info() export
		// 3. Check for protocol conflicts with already-installed plugins
		// 4. Store in ~/.apiqube/plugins/<name>.wasm
		// 5. Print plugin info on success
		return fmt.Errorf("not implemented")
	},
}

var pluginListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// 1. Walk plugin directory
		// 2. For each .wasm file, load metadata via plugin_info()
		// 3. Print table: name, version, protocols, fields count
		return fmt.Errorf("not implemented")
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an installed plugin",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		return fmt.Errorf("not implemented")
	},
}

func init() {
	pluginCmd.AddCommand(pluginInstallCmd, pluginListCmd, pluginRemoveCmd)
}
