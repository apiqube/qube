package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/apiqube/engine"

	"github.com/apiqube/qube/internal/discovery"
	"github.com/apiqube/qube/internal/output"
	"github.com/apiqube/qube/internal/ui"
	uirun "github.com/apiqube/qube/internal/ui/run"
)

var runFlags struct {
	configPath   string
	pluginDir    string
	tags         []string
	excludeTags  []string
	parallel     bool
	failFast     bool
	timeout      string
	outputFormat string
	verbose      bool
	envFile      string
	noTUI        bool
}

var runCmd = &cobra.Command{
	Use:   "run [path...]",
	Short: "Run tests from files or directories",
	Long: `Run executes tests from the given YAML files and directories.

If no path is given, qube looks for tests in the ./tests/ directory.
Uses .qube.yaml from the current or ancestor directories if present.`,
	Args: cobra.ArbitraryArgs,
	RunE: runE,
}

func runE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cwd, _ := os.Getwd()
	configPath := runFlags.configPath
	if configPath == "" {
		configPath = discovery.FindConfig(cwd)
	}

	envPath := runFlags.envFile
	if envPath == "" {
		envPath = discovery.FindEnvFile(cwd)
	}
	envVars, err := discovery.LoadEnvFile(envPath)
	if err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	paths := args
	if len(paths) == 0 {
		paths = []string{"./tests/"}
	}

	warnUnsupportedFlags(cmd.ErrOrStderr())

	eng := engine.New(
		engine.WithParallel(runFlags.parallel),
		engine.WithFailFast(runFlags.failFast),
		engine.WithPluginDir(resolvePluginDir(runFlags.pluginDir)),
	)
	defer eng.Close()

	var runOpts []engine.RunOption
	if envVars != nil {
		runOpts = append(runOpts, engine.WithEnv(envVars))
	}
	if configPath != "" {
		runOpts = append(runOpts, engine.WithConfigPath(configPath))
	}

	results, runErr := dispatchRun(ctx, cmd, eng, engine.FromPaths(paths...), runOpts)

	if runErr != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), ui.Failure.Render("engine error: ")+runErr.Error())
		os.Exit(2) //nolint:gocritic // explicit exit code per CLI convention
	}
	if results != nil && (results.Failed > 0 || results.Errored > 0) {
		os.Exit(1) //nolint:gocritic
	}
	return nil
}

// dispatchRun picks the right output sink based on flags and TTY detection,
// then drives engine.Run.
func dispatchRun(
	ctx context.Context,
	cmd *cobra.Command,
	eng *engine.Engine,
	input engine.Input,
	opts []engine.RunOption,
) (*engine.Results, error) {
	stdout := cmd.OutOrStdout()

	switch runFlags.outputFormat {
	case "json":
		h := output.NewJSON(stdout)
		opts = append(opts, engine.WithHandler(h))
		return eng.Run(ctx, input, opts...)

	case "junit":
		h := output.NewJUnit(stdout)
		opts = append(opts, engine.WithHandler(h))
		return eng.Run(ctx, input, opts...)

	case "tap":
		h := output.NewTAP(stdout)
		opts = append(opts, engine.WithHandler(h))
		return eng.Run(ctx, input, opts...)

	default:
		// pretty: TUI when interactive, progressive lipgloss otherwise
		if !runFlags.noTUI && ui.IsInteractive() {
			return uirun.Run(ctx, eng, input, opts...)
		}
		h := output.NewPretty(stdout, runFlags.verbose)
		opts = append(opts, engine.WithHandler(h))
		return eng.Run(ctx, input, opts...)
	}
}

func warnUnsupportedFlags(w io.Writer) {
	if len(runFlags.tags) > 0 || len(runFlags.excludeTags) > 0 {
		msg := "tag filtering not yet supported in v1.0; flags accepted, no effect"
		fmt.Fprintln(w, ui.Warn.Render("⚠ ")+ui.Muted.Render(msg))
	}
	if runFlags.timeout != "" {
		msg := "global timeout flag not yet wired; ignored in v1.0"
		fmt.Fprintln(w, ui.Warn.Render("⚠ ")+ui.Muted.Render(msg))
	}
}

func init() {
	runCmd.Flags().StringVarP(&runFlags.configPath, "config", "c", "", "path to .qube.yaml")
	runCmd.Flags().StringVar(&runFlags.pluginDir, "plugins", "", "plugin directory")
	runCmd.Flags().StringSliceVarP(&runFlags.tags, "tags", "t", nil, "only run tests with these tags (v1.0: accepted, no effect)")
	runCmd.Flags().StringSliceVar(&runFlags.excludeTags, "exclude-tags", nil, "skip tests with these tags (v1.0: accepted, no effect)")
	runCmd.Flags().BoolVar(&runFlags.parallel, "parallel", true, "enable parallel execution")
	runCmd.Flags().BoolVar(&runFlags.failFast, "fail-fast", false, "stop on first failure")
	runCmd.Flags().StringVar(&runFlags.timeout, "timeout", "", "global timeout override (v1.0: accepted, no effect)")
	runCmd.Flags().StringVarP(&runFlags.outputFormat, "output", "o", "pretty", "output format: pretty, json, junit, tap")
	runCmd.Flags().BoolVarP(&runFlags.verbose, "verbose", "v", false, "detailed output")
	runCmd.Flags().StringVar(&runFlags.envFile, "env-file", "", "load environment variables from file")
	runCmd.Flags().BoolVar(&runFlags.noTUI, "no-tui", false, "force progressive output even on interactive terminals")
}
