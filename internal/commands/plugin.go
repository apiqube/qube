package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		ui.Warn("no plugin directory configured")
		ui.Info("pass --plugins, set $QUBE_PLUGIN_DIR, or create ~/.apiqube/plugins/")
		return nil
	}

	entries, err := walkPluginDir(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		ui.Infof("no plugins installed in %s", ui.AccentStyle.Render(dir))
		ui.Info("add a *.wasm file there or run 'qube plugin install <name>' (coming soon)")
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

	// Trigger lazy plugin load via an empty Run so we can snapshot Info.
	_, _ = eng.Run(ctx, engine.FromBytes([]byte("tests: []")))
	return eng.Plugins(), nil
}

// renderPluginTable returns a no-frills aligned table: bold header, separator
// rule, plain rows. No borders, no zebra-stripes — just structured columns.
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
	widths := columnWidths(headers, rows)

	var b strings.Builder
	b.WriteString(formatRow(headers, widths, func(s string) string { return ui.TableHeaderStyle.Render(s) }))
	b.WriteByte('\n')
	b.WriteString(ui.MutedStyle.Render(strings.Repeat("─", totalWidth(widths))))
	b.WriteByte('\n')
	for _, row := range rows {
		b.WriteString(formatRow(row, widths, func(s string) string { return s }))
		b.WriteByte('\n')
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

func columnWidths(headers []string, rows [][]string) []int {
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

func totalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w
	}
	return total + (len(widths)-1)*2 // join separator "  "
}

func formatRow(cells []string, widths []int, render func(string) string) string {
	parts := make([]string, len(cells))
	for i, c := range cells {
		pad := widths[i] - len(c)
		if pad < 0 {
			pad = 0
		}
		parts[i] = render(c + strings.Repeat(" ", pad))
	}
	return strings.Join(parts, "  ")
}

func printRoadmapStub(stdout io.Writer, command, reason string) {
	fmt.Fprintln(stdout, ui.AccentStyle.Render(command)+ui.MutedStyle.Render(" — not yet implemented"))
	fmt.Fprintln(stdout, ui.MutedStyle.Render("  "+reason))
	fmt.Fprintln(stdout, ui.MutedStyle.Render("  v1.0 covers: ")+ui.AccentStyle.Render("run, check, init, version, plugin list"))
}

func init() {
	pluginListCmd.Flags().StringVar(&pluginListFlags.dir, "plugins", "",
		"plugin directory (defaults to $QUBE_PLUGIN_DIR or ~/.apiqube/plugins)")
	pluginCmd.AddCommand(pluginInstallCmd, pluginListCmd, pluginRemoveCmd)
}
