package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/apiqube/qube/internal/ui"
	"github.com/apiqube/qube/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print qube version",
	RunE: func(cmd *cobra.Command, _ []string) error {
		out := cmd.OutOrStdout()
		title := ui.BrandStyle.Render("qube ") + ui.AccentStyle.Render(version.Version)
		meta := ui.MutedStyle.Render(fmt.Sprintf(" · commit %s · built %s",
			version.Commit, version.Date))
		fmt.Fprintln(out, title+meta)
		return nil
	},
}
