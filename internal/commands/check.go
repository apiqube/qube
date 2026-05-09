package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

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
		fmt.Fprintln(cmd.ErrOrStderr(), ui.Failure.Render("check failed: ")+err.Error())
		os.Exit(2) //nolint:gocritic
	}

	stdout := cmd.OutOrStdout()
	if len(errs) == 0 {
		fmt.Fprintln(stdout, ui.SummaryCard.
			BorderForeground(ui.ColorSuccess).
			Render(ui.Success.Render("All manifests valid.")))
		return nil
	}

	renderErrors(stdout, errs)
	os.Exit(1) //nolint:gocritic
	return nil
}

func renderErrors(stdout io.Writer, errs []engine.ValidationError) {
	byFile := groupErrorsByFile(errs)
	files := make([]string, 0, len(byFile))
	for f := range byFile {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		header := ui.Brand.Render(file)
		var lines []string
		for _, e := range byFile[file] {
			icon := ui.Failure.Render("✗")
			if e.Severity == engine.SeverityWarning {
				icon = ui.Warn.Render("⚠")
			}
			loc := ""
			if e.Line > 0 {
				loc = ui.Muted.Render(fmt.Sprintf("line %d  ", e.Line))
			}
			field := ""
			if e.Field != "" {
				field = ui.Accent.Render(e.Field) + "  "
			}
			lines = append(lines, fmt.Sprintf("  %s  %s%s%s", icon, loc, field, e.Message))
		}
		body := header + "\n" + strings.Join(lines, "\n")
		fmt.Fprintln(stdout, ui.ErrorBlock.Render(body))
	}

	total := ui.Failure.Render(fmt.Sprintf("%d validation error(s) found.", len(errs)))
	fmt.Fprintln(stdout, total)
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
