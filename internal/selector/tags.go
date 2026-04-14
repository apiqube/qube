package selector

// MatchTags reports whether a test with the given tags should be included
// based on the user's --tags and --exclude-tags flags.
//
// Rules:
//   - If includeTags is empty, all tests are included by default
//   - If includeTags is non-empty, test must have at least one matching tag
//   - If excludeTags is non-empty and test has any matching tag, it is excluded
//   - Exclusion always wins over inclusion
func MatchTags(testTags, includeTags, excludeTags []string) bool {
	// TODO: implementation
	//
	// 1. Build set from excludeTags
	// 2. If any testTag in excludeSet → return false (skip)
	// 3. If includeTags is empty → return true (all tests)
	// 4. Build set from includeTags
	// 5. If any testTag in includeSet → return true
	// 6. Return false (no matching tag)
	return true
}
