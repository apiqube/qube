package commands

import (
	"fmt"

	"github.com/spf13/cobra"
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
}

var runCmd = &cobra.Command{
	Use:   "run [path...]",
	Short: "Run tests from files or directories",
	Long: `Run executes tests from the given YAML files and directories.

If no path is given, qube looks for tests in the ./tests/ directory.
Uses .qube.yaml from the current or ancestor directories if present.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementation
		//
		// 1. Discover config: .qube.yaml walking up from cwd (unless --config)
		// 2. Load .env file if present or specified via --env-file
		// 3. Resolve paths (default to ./tests/ if args empty)
		// 4. Build engine.Engine with options from config + CLI flags
		// 5. Select output handler based on --output flag:
		//    - pretty (default) → output.NewPretty()
		//    - json             → output.NewJSON(stdout)
		//    - junit            → output.NewJUnit(stdout)
		// 6. Apply tag filters via selector.Filter(tests, includeTags, excludeTags)
		// 7. Call eng.Run(ctx, engine.FromPaths(paths...), options...)
		// 8. Exit with non-zero code if any test failed
		return fmt.Errorf("not implemented")
	},
}

func init() {
	runCmd.Flags().StringVarP(&runFlags.configPath, "config", "c", "", "path to .qube.yaml")
	runCmd.Flags().StringVar(&runFlags.pluginDir, "plugins", "", "plugin directory")
	runCmd.Flags().StringSliceVarP(&runFlags.tags, "tags", "t", nil, "only run tests with these tags")
	runCmd.Flags().StringSliceVar(&runFlags.excludeTags, "exclude-tags", nil, "skip tests with these tags")
	runCmd.Flags().BoolVar(&runFlags.parallel, "parallel", true, "enable parallel execution")
	runCmd.Flags().BoolVar(&runFlags.failFast, "fail-fast", false, "stop on first failure")
	runCmd.Flags().StringVar(&runFlags.timeout, "timeout", "", "global timeout override")
	runCmd.Flags().StringVarP(&runFlags.outputFormat, "output", "o", "pretty", "output format: pretty, json, junit, tap")
	runCmd.Flags().BoolVarP(&runFlags.verbose, "verbose", "v", false, "detailed output")
	runCmd.Flags().StringVar(&runFlags.envFile, "env-file", "", "load environment variables from file")
}
