package discovery

// FindConfig walks upward from the given directory looking for a .qube.yaml file.
// Returns the absolute path if found, or empty string if not.
func FindConfig(startDir string) string {
	// TODO: implementation
	//
	// 1. Start at startDir
	// 2. Look for .qube.yaml in current directory
	// 3. If found, return absolute path
	// 4. If not, move to parent directory
	// 5. Stop at filesystem root or at $HOME (whichever comes first)
	return ""
}
