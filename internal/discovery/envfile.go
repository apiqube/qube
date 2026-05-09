package discovery

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// EnvFileName is the .env file discovery walks toward.
const EnvFileName = ".env"

// FindEnvFile returns the absolute path to a .env in startDir or one of its
// ancestors (up to $HOME or filesystem root), or "" if not found.
func FindEnvFile(startDir string) string {
	return findUpward(startDir, EnvFileName)
}

// LoadEnvFile parses a .env file and returns its key-value pairs.
//
// Supported syntax:
//
//	KEY=value             # plain
//	KEY="value"           # double-quoted, supports escaped \" \\ \n \t
//	KEY='value'           # single-quoted, no interpretation
//	# comment             # whole-line comment
//	KEY=value # trailing  # trailing comment after unquoted value
//	export KEY=value      # leading "export " is tolerated
//
// Returns (nil, nil) if path is empty or the file does not exist — both are
// non-error states (no .env in cwd is the common case).
func LoadEnvFile(path string) (map[string]string, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	out := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return nil, fmt.Errorf("%s:%d: missing '='", path, lineNo)
		}
		key := strings.TrimSpace(line[:eq])
		if key == "" {
			return nil, fmt.Errorf("%s:%d: empty key", path, lineNo)
		}
		value, err := parseValue(strings.TrimSpace(line[eq+1:]))
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNo, err)
		}
		out[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// parseValue handles unquoted (with optional trailing comment), single-quoted,
// and double-quoted (with escapes) value forms.
func parseValue(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	switch raw[0] {
	case '"':
		return parseDoubleQuoted(raw)
	case '\'':
		return parseSingleQuoted(raw)
	}

	// Unquoted: strip trailing " # comment".
	if idx := strings.Index(raw, " #"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	return raw, nil
}

func parseDoubleQuoted(raw string) (string, error) {
	if len(raw) < 2 || raw[len(raw)-1] != '"' {
		// Allow trailing comment after the closing quote.
		end := strings.LastIndexByte(raw, '"')
		if end <= 0 {
			return "", errors.New(`unterminated double-quoted value`)
		}
		raw = raw[:end+1]
	}
	body := raw[1 : len(raw)-1]
	var b strings.Builder
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c != '\\' || i == len(body)-1 {
			b.WriteByte(c)
			continue
		}
		i++
		switch body[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '\\':
			b.WriteByte('\\')
		case '"':
			b.WriteByte('"')
		default:
			b.WriteByte('\\')
			b.WriteByte(body[i])
		}
	}
	return b.String(), nil
}

func parseSingleQuoted(raw string) (string, error) {
	if len(raw) < 2 || raw[len(raw)-1] != '\'' {
		end := strings.LastIndexByte(raw, '\'')
		if end <= 0 {
			return "", errors.New(`unterminated single-quoted value`)
		}
		raw = raw[:end+1]
	}
	return raw[1 : len(raw)-1], nil
}
