package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/ui"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage qube plugins",
	Long:  `Plugin manages WASM plugins that extend qube with new protocols, reporters, generators, and hooks.`,
}

var pluginListFlags struct {
	dir string
}

var pluginListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed plugins",
	RunE:    pluginListE,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name|url|path>",
	Short: "Install a plugin (not yet implemented)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, _ []string) error {
		printRoadmapStub(cmd.OutOrStdout(), "qube plugin install",
			"installing from a name/url/path requires a plugin registry and version pinning.")
		return nil
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an installed plugin (not yet implemented)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, _ []string) error {
		printRoadmapStub(cmd.OutOrStdout(), "qube plugin remove",
			"depends on the registry/install pipeline that's still being designed.")
		return nil
	},
}

func pluginListE(cmd *cobra.Command, _ []string) error {
	dir := resolvePluginDir(pluginListFlags.dir)
	if dir == "" {
		fmt.Fprintln(cmd.OutOrStdout(), ui.SummaryCard.Render(
			ui.Muted.Render("No plugin directory configured.\n")+
				ui.Muted.Render("Pass --plugins, set $QUBE_PLUGIN_DIR, or create ~/.apiqube/plugins/")))
		return nil
	}

	entries, err := walkPluginDir(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), ui.SummaryCard.Render(
			ui.Muted.Render("No plugins installed in ")+ui.Accent.Render(dir)+".\n"+
				ui.Muted.Render("Add a *.wasm file there or run ")+ui.Accent.Render("qube plugin install <name>")+ui.Muted.Render(" (coming soon).")))
		return nil
	}

	plugins, err := loadPluginInfo(cmd.Context(), dir)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), renderPluginTable(plugins))
	return nil
}

func resolvePluginDir(flag string) string {
	if flag != "" {
		return flag
	}
	if v := os.Getenv("QUBE_PLUGIN_DIR"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	candidate := filepath.Join(home, ".apiqube", "plugins")
	if _, err := os.Stat(candidate); err != nil {
		return ""
	}
	return candidate
}

func walkPluginDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read plugin dir %s: %w", dir, err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".wasm") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out, nil
}

func loadPluginInfo(ctx context.Context, dir string) ([]engine.PluginSchema, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	eng := engine.New(engine.WithPluginDir(dir))
	defer eng.Close()

	// Trigger lazy plugin load by performing a benign Run that does nothing.
	// We pipe an empty manifest so engine instantiates plugins and we can
	// then snapshot them.
	_, _ = eng.Run(ctx, engine.FromBytes([]byte("tests: []")))
	return eng.Plugins(), nil
}

func renderPluginTable(plugins []engine.PluginSchema) string {
	headers := []string{"NAME", "VERSION", "PROTOCOLS", "CAPABILITIES"}
	rows := make([][]string, 0, len(plugins))
	for _, p := range plugins {
		rows = append(rows, []string{
			p.Name,
			p.Version,
			joinProtocols(p.Protocols),
			strings.Join(p.Capabilities, ", "),
		})
	}

	widths := computeColumnWidths(headers, rows)
	var b strings.Builder
	b.WriteString(renderTableRow(headers, widths, ui.TableHeader))
	b.WriteString("\n")
	for i, row := range rows {
		style := ui.TableCell
		if i%2 == 1 {
			style = ui.TableCellFaint
		}
		b.WriteString(renderTableRow(row, widths, style))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func joinProtocols(p []engine.Protocol) string {
	out := make([]string, len(p))
	for i, x := range p {
		out[i] = string(x)
	}
	return strings.Join(out, ", ")
}

func computeColumnWidths(headers []string, rows [][]string) []int {
	w := make([]int, len(headers))
	for i, h := range headers {
		w[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > w[i] {
				w[i] = len(cell)
			}
		}
	}
	return w
}

func renderTableRow(cells []string, widths []int, style lipgloss.Style) string {
	parts := make([]string, len(cells))
	for i, c := range cells {
		padding := widths[i] - len(c)
		if padding < 0 {
			padding = 0
		}
		parts[i] = style.Render(c + strings.Repeat(" ", padding))
	}
	return strings.Join(parts, "  ")
}

func printRoadmapStub(stdout interface{ Write(p []byte) (n int, err error) }, command, reason string) {
	body := strings.Join([]string{
		ui.Brand.Render(command),
		"",
		ui.Muted.Render("Roadmap status: not yet implemented."),
		ui.Muted.Render(reason),
		"",
		ui.Muted.Render("v1.0 covers: ") + ui.Accent.Render("run, check, init, version, plugin list"),
		ui.Muted.Render("Track progress at https://github.com/apiqube/qube"),
	}, "\n")
	fmt.Fprintln(stdout, ui.Card.Render(body))
}

func init() {
	pluginListCmd.Flags().StringVar(&pluginListFlags.dir, "plugins", "",
		"plugin directory (defaults to $QUBE_PLUGIN_DIR or ~/.apiqube/plugins)")
	pluginCmd.AddCommand(pluginInstallCmd, pluginListCmd, pluginRemoveCmd)
}
