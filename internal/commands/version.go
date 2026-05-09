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
		title := ui.Brand.Render("qube ") + ui.Accent.Render(version.Version)
		details := ui.Muted.Render(fmt.Sprintf("commit %s · built %s",
			version.Commit, version.Date))
		fmt.Fprintln(out, ui.SummaryCard.Render(title+"\n"+details))
		return nil
	},
}
