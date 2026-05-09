package tui

import "testing"

func TestValidateURL(t *testing.T) {
	cases := map[string]bool{
		"":                     false, // required
		"  ":                   false, // whitespace
		"localhost:8080":       false, // missing scheme
		"http://localhost":     true,
		"https://api.test/v1":  true,
	}
	for in, ok := range cases {
		err := validateURL(in)
		if (err == nil) != ok {
			t.Errorf("validateURL(%q) ok=%v; want %v (err=%v)", in, err == nil, ok, err)
		}
	}
}

func TestInitAnswers_ZeroValue(t *testing.T) {
	var a InitAnswers
	if a.TargetURL != "" || a.HasSwagger || a.SwaggerURL != "" || len(a.Plugins) != 0 {
		t.Errorf("zero-value should be empty: %+v", a)
	}
}
