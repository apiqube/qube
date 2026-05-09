package selector

import "testing"

func TestMatchTags(t *testing.T) {
	cases := []struct {
		name        string
		testTags    []string
		include     []string
		exclude     []string
		want        bool
	}{
		{"no filters → include", []string{"a"}, nil, nil, true},
		{"include miss", []string{"a"}, []string{"b"}, nil, false},
		{"include hit", []string{"a", "b"}, []string{"b"}, nil, true},
		{"exclude wins", []string{"a", "b"}, nil, []string{"a"}, false},
		{"exclude wins over include", []string{"a"}, []string{"a"}, []string{"a"}, false},
		{"include with empty test tags", nil, []string{"a"}, nil, false},
		{"exclude with empty test tags", nil, nil, []string{"a"}, true},
		{"both empty test tags", nil, nil, nil, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := MatchTags(c.testTags, c.include, c.exclude)
			if got != c.want {
				t.Errorf("MatchTags(%v, %v, %v) = %v; want %v",
					c.testTags, c.include, c.exclude, got, c.want)
			}
		})
	}
}
