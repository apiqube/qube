package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoot_HasAllSubcommands(t *testing.T) {
	r := Root()
	want := []string{"run", "check", "init", "generate", "plugin", "version"}
	for _, name := range want {
		_, _, err := r.Find([]string{name})
		if err != nil {
			t.Errorf("subcommand %q missing: %v", name, err)
		}
	}
}

func TestVersionCommand_Output(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer
	cmd := versionCmd
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("version RunE: %v", err)
	}
	if !strings.Contains(out.String(), "qube") {
		t.Errorf("version output missing brand: %q", out.String())
	}
}

func TestGenerateCommand_StubMessage(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer
	cmd := generateCmd
	cmd.SetOut(&out)
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "qube generate") {
		t.Errorf("stub output missing command name: %q", out.String())
	}
	if !strings.Contains(out.String(), "not yet implemented") {
		t.Errorf("stub output missing roadmap text: %q", out.String())
	}
}

func TestPluginInstall_StubMessage(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer
	cmd := pluginInstallCmd
	cmd.SetOut(&out)
	if err := cmd.RunE(cmd, []string{"some-plugin"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "not yet implemented") {
		t.Errorf("install stub message missing: %q", out.String())
	}
}

func TestPluginRemove_StubMessage(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer
	cmd := pluginRemoveCmd
	cmd.SetOut(&out)
	if err := cmd.RunE(cmd, []string{"some-plugin"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "not yet implemented") {
		t.Errorf("remove stub message missing: %q", out.String())
	}
}

func TestInit_SilentMode(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	dir := t.TempDir()
	prevWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWD)

	initFlags.interactive = false
	initFlags.targetURL = "http://example.test"
	initFlags.force = false
	defer func() {
		initFlags.targetURL = ""
		initFlags.force = false
	}()

	var out bytes.Buffer
	cmd := initCmd
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("init silent: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".qube.yaml")); err != nil {
		t.Errorf(".qube.yaml not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tests", "example.yaml")); err != nil {
		t.Errorf("tests/example.yaml not written: %v", err)
	}
	if !strings.Contains(out.String(), "Project initialized") {
		t.Errorf("success message missing: %q", out.String())
	}
}

func TestInit_RefusesOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	prevWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWD)

	// Pre-existing .qube.yaml
	if err := os.WriteFile(filepath.Join(dir, ".qube.yaml"), []byte("v: 1"), 0o644); err != nil {
		t.Fatal(err)
	}

	initFlags.interactive = false
	initFlags.force = false
	defer func() { initFlags.force = false }()

	var out bytes.Buffer
	cmd := initCmd
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err == nil {
		t.Error("init should refuse to overwrite without --force")
	}
}

func TestRenderPluginTable_NotEmpty(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	plugins := mockPluginSchemas()
	rendered := renderPluginTable(plugins)
	if !strings.Contains(rendered, "NAME") {
		t.Error("table missing header")
	}
	if !strings.Contains(rendered, "http") {
		t.Error("table missing http row")
	}
}

func TestResolvePluginDir_FlagWins(t *testing.T) {
	t.Setenv("QUBE_PLUGIN_DIR", "/from/env")
	got := resolvePluginDir("/from/flag")
	if got != "/from/flag" {
		t.Errorf("flag should win, got %q", got)
	}
}

func TestResolvePluginDir_EnvFallback(t *testing.T) {
	t.Setenv("QUBE_PLUGIN_DIR", "/from/env")
	got := resolvePluginDir("")
	if got != "/from/env" {
		t.Errorf("env should be used, got %q", got)
	}
}
