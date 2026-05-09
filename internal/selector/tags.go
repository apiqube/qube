package selector

// MatchTags reports whether a test with the given tags should be included
// based on the user's --tags and --exclude-tags flags.
//
// Rules:
//   - excludeTags wins: any matching tag → excluded.
//   - includeTags is a positive filter: empty means "all"; non-empty means
//     the test must have at least one matching tag.
func MatchTags(testTags, includeTags, excludeTags []string) bool {
	if anyMatch(testTags, excludeTags) {
		return false
	}
	if len(includeTags) == 0 {
		return true
	}
	return anyMatch(testTags, includeTags)
}

func anyMatch(have, want []string) bool {
	if len(have) == 0 || len(want) == 0 {
		return false
	}
	wantSet := make(map[string]struct{}, len(want))
	for _, t := range want {
		wantSet[t] = struct{}{}
	}
	for _, t := range have {
		if _, ok := wantSet[t]; ok {
			return true
		}
	}
	return false
}
