package commands

import (
	"fmt"

	"github.com/apiqube/qube/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print qube version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "qube %s (%s) built %s\n",
			version.Version, version.Commit, version.Date)
		return nil
	},
}
