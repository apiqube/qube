package discovery

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadEnvFile_EmptyPath(t *testing.T) {
	got, err := LoadEnvFile("")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("empty path should yield nil, got %v", got)
	}
}

func TestLoadEnvFile_Missing(t *testing.T) {
	got, err := LoadEnvFile(filepath.Join(t.TempDir(), "nope.env"))
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("missing file should yield nil, got %v", got)
	}
}

func TestLoadEnvFile_PlainPairs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	contents := `
# comment line
KEY1=value1
KEY2=value 2
export EXPORTED=ok

KEY3 = trimmed
TRAILING=value # trailing comment
`
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"KEY1":     "value1",
		"KEY2":     "value 2",
		"EXPORTED": "ok",
		"KEY3":     "trimmed",
		"TRAILING": "value",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestLoadEnvFile_DoubleQuotedEscapes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	contents := `WITH_NL="line1\nline2"
WITH_TAB="a\tb"
WITH_QUOTE="say \"hi\""
WITH_SLASH="a\\b"
`
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		"WITH_NL":    "line1\nline2",
		"WITH_TAB":   "a\tb",
		"WITH_QUOTE": `say "hi"`,
		"WITH_SLASH": `a\b`,
	}
	for k, want := range cases {
		if got[k] != want {
			t.Errorf("%s = %q; want %q", k, got[k], want)
		}
	}
}

func TestLoadEnvFile_SingleQuoted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	contents := "RAW='line1\\nline2 # not a comment'\n"
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := LoadEnvFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["RAW"] != `line1\nline2 # not a comment` {
		t.Errorf("single-quoted should not interpret escapes/comments: %q", got["RAW"])
	}
}

func TestLoadEnvFile_MissingEquals(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("BADLINE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadEnvFile(path); err == nil {
		t.Error("missing '=' should error")
	}
}

func TestLoadEnvFile_EmptyKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("=value\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadEnvFile(path); err == nil {
		t.Error("empty key should error")
	}
}

func TestFindEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("X=Y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindEnvFile(dir)
	want, _ := filepath.Abs(envPath)
	if got != want {
		t.Errorf("FindEnvFile = %q; want %q", got, want)
	}
}
