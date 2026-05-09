package commands

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/apiqube/qube/internal/tui"
	"github.com/apiqube/qube/internal/ui"
)

var initFlags struct {
	interactive bool
	swaggerURL  string
	targetURL   string
	force       bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new qube project",
	Long: `Init creates starter files for a new qube project:
  .qube.yaml         project configuration
  tests/example.yaml simple example test

Use --interactive to answer questions and (later) generate tests from
an existing OpenAPI/Swagger specification.`,
	RunE: initE,
}

func initE(cmd *cobra.Command, _ []string) error {
	target := initFlags.targetURL
	if target == "" {
		target = "http://localhost:8080"
	}

	if initFlags.interactive {
		answers, err := tui.RunInitWizard()
		if err != nil {
			return fmt.Errorf("init wizard: %w", err)
		}
		if answers != nil {
			target = answers.TargetURL
			initFlags.swaggerURL = answers.SwaggerURL
		}
	}

	if err := writeProjectFiles(cmd.OutOrStdout(), target); err != nil {
		return err
	}
	return nil
}

// writeProjectFiles writes .qube.yaml and tests/example.yaml relative to cwd.
// Honors --force; refuses to overwrite otherwise. Output is structured log
// lines, no boxed cards.
func writeProjectFiles(_ io.Writer, target string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(cwd, ".qube.yaml")
	if err := writeIfAbsent(cfgPath, configYAML(target), initFlags.force); err != nil {
		return err
	}

	exDir := filepath.Join(cwd, "tests")
	if err := os.MkdirAll(exDir, 0o755); err != nil {
		return fmt.Errorf("create tests dir: %w", err)
	}
	exPath := filepath.Join(exDir, "example.yaml")
	if err := writeIfAbsent(exPath, exampleYAML(target), initFlags.force); err != nil {
		return err
	}

	ui.Success("Project initialized.")
	ui.Infof("Created %s", cfgPath)
	ui.Infof("Created %s", exPath)
	ui.Infof("Next: %s", ui.AccentStyle.Render("qube run tests/"))
	return nil
}

func writeIfAbsent(path, contents string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s exists; pass --force to overwrite", path)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func configYAML(target string) string {
	return fmt.Sprintf(`# qube project configuration. See https://github.com/apiqube/qube
version: 1
targets:
  default: %s
runner:
  parallel: true
  failFast: false
`, target)
}

func exampleYAML(target string) string {
	return fmt.Sprintf(`# Example test. Run with: qube run tests/example.yaml
target: %s

tests:
  - name: Health check
    method: GET
    resource: /health
    expect:
      status: 200
`, target)
}

func init() {
	initCmd.Flags().BoolVarP(&initFlags.interactive, "interactive", "i", false, "ask questions interactively")
	initCmd.Flags().StringVar(&initFlags.swaggerURL, "from-swagger", "", "generate tests from OpenAPI/Swagger spec (v1.0: not yet wired)")
	initCmd.Flags().StringVar(&initFlags.targetURL, "target", "", "default target URL")
	initCmd.Flags().BoolVar(&initFlags.force, "force", false, "overwrite existing files")
}
