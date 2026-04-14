package discovery

// LoadEnvFile parses a .env file and returns its key-value pairs.
// Supports simple KEY=VALUE lines, comments starting with #, and quoted values.
func LoadEnvFile(path string) (map[string]string, error) {
	// TODO: implementation
	//
	// 1. Open file (return nil, nil if not found — missing .env is not an error)
	// 2. Read line by line:
	//    - Skip empty lines and lines starting with #
	//    - Parse KEY=VALUE
	//    - Handle quoted values (VALUE="with spaces")
	//    - Handle escaped characters in quoted values
	// 3. Return map
	return nil, nil
}

// FindEnvFile returns the path to .env in startDir or its parents, or empty if not found.
func FindEnvFile(startDir string) string {
	// TODO: implementation (similar to FindConfig but looking for .env)
	return ""
}
