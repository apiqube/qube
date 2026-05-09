package commands

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/discovery"
	"github.com/apiqube/qube/internal/ui"
)

var checkFlags struct {
	configPath string
	pluginDir  string
}

var checkCmd = &cobra.Command{
	Use:   "check [path...]",
	Short: "Validate test manifests without executing them",
	Long: `Check validates YAML syntax, plugin requirements, unresolved template
references, and dependency cycles.

Exits with non-zero status if any errors are found.`,
	Args: cobra.ArbitraryArgs,
	RunE: checkE,
}

func checkE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cwd, _ := os.Getwd()
	configPath := checkFlags.configPath
	if configPath == "" {
		configPath = discovery.FindConfig(cwd)
	}

	paths := args
	if len(paths) == 0 {
		paths = []string{"./tests/"}
	}

	eng := engine.New(engine.WithPluginDir(resolvePluginDir(checkFlags.pluginDir)))
	defer eng.Close()

	var opts []engine.CheckOption
	if configPath != "" {
		opts = append(opts, engine.WithCheckConfigPath(configPath))
	}

	errs, err := eng.Check(ctx, engine.FromPaths(paths...), opts...)
	if err != nil {
		ui.Errf("check failed: %v", err)
		os.Exit(2) //nolint:gocritic
	}

	if len(errs) == 0 {
		ui.Success("All manifests valid")
		return nil
	}

	renderErrors(errs)
	os.Exit(1) //nolint:gocritic
	return nil
}

func renderErrors(errs []engine.ValidationError) {
	byFile := groupErrorsByFile(errs)
	files := make([]string, 0, len(byFile))
	for f := range byFile {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		ui.Err(file)
		for _, e := range byFile[file] {
			icon := "✗"
			level := ui.LevelError
			if e.Severity == engine.SeverityWarning {
				icon = "⚠"
				level = ui.LevelWarn
			}
			detail := formatErrorDetail(icon, e)
			ui.DefaultLogger.Log(level, detail)
		}
	}
	ui.Errf("%d validation error(s) found", len(errs))
}

func formatErrorDetail(icon string, e engine.ValidationError) string {
	parts := []string{icon}
	if e.Line > 0 {
		parts = append(parts, ui.MutedStyle.Render(fmt.Sprintf("line %d", e.Line)))
	}
	if e.Field != "" {
		parts = append(parts, ui.AccentStyle.Render(e.Field))
	}
	parts = append(parts, e.Message)
	return joinParts(parts, "  ")
}

func joinParts(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

func groupErrorsByFile(errs []engine.ValidationError) map[string][]engine.ValidationError {
	out := map[string][]engine.ValidationError{}
	for _, e := range errs {
		key := e.File
		if key == "" {
			key = "<input>"
		}
		out[key] = append(out[key], e)
	}
	return out
}

func init() {
	checkCmd.Flags().StringVarP(&checkFlags.configPath, "config", "c", "", "path to .qube.yaml")
	checkCmd.Flags().StringVar(&checkFlags.pluginDir, "plugins", "", "plugin directory")
}
