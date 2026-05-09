package discovery

import (
	"os"
	"path/filepath"
)

// ConfigFileName is the file name discovery walks toward.
const ConfigFileName = ".qube.yaml"

// FindConfig walks upward from startDir looking for a .qube.yaml file. It
// stops at the user's home directory or the filesystem root, whichever
// comes first. Returns the absolute path or "" if not found.
func FindConfig(startDir string) string {
	return findUpward(startDir, ConfigFileName)
}

// findUpward is the shared walk-up helper used by both FindConfig and
// FindEnvFile. It resolves startDir to an absolute path and ascends until
// it hits $HOME, the filesystem root, or finds the target file.
func findUpward(startDir, fileName string) string {
	abs, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}

	home, _ := os.UserHomeDir()

	for {
		candidate := filepath.Join(abs, fileName)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}

		parent := filepath.Dir(abs)
		if parent == abs {
			return "" // hit the root
		}
		if home != "" && abs == home {
			return "" // do not search above HOME
		}
		abs = parent
	}
}
