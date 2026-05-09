package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfig_InCurrentDir(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, ".qube.yaml")
	if err := os.WriteFile(cfg, []byte("version: 1"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindConfig(dir)
	want, _ := filepath.Abs(cfg)
	if got != want {
		t.Errorf("FindConfig = %q; want %q", got, want)
	}
}

func TestFindConfig_InAncestor(t *testing.T) {
	root := t.TempDir()
	mid := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(mid, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(root, ".qube.yaml")
	if err := os.WriteFile(cfg, []byte("version: 1"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindConfig(mid)
	want, _ := filepath.Abs(cfg)
	if got != want {
		t.Errorf("FindConfig from descendant = %q; want %q", got, want)
	}
}

func TestFindConfig_NotFound(t *testing.T) {
	dir := t.TempDir()
	if got := FindConfig(dir); got != "" {
		t.Errorf("expected empty result, got %q", got)
	}
}

func TestFindConfig_StopsAtHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // Windows

	// Place config ABOVE the home dir; FindConfig should not climb past home.
	above := filepath.Dir(home)
	cfg := filepath.Join(above, ".qube.yaml")
	if err := os.WriteFile(cfg, []byte("v: 1"), 0o644); err != nil {
		t.Skipf("cannot write above tempdir: %v", err)
	}
	defer os.Remove(cfg)

	sub := filepath.Join(home, "project")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	got := FindConfig(sub)
	if got != "" {
		t.Errorf("FindConfig should not climb above HOME; got %q", got)
	}
}

func TestFindConfig_BadStartDir(t *testing.T) {
	// FindConfig handles its own absolution; an empty string maps to cwd.
	// The test simply asserts no panic and a string result.
	_ = FindConfig("")
}
