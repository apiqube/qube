package ui

import "testing"

func TestStatusStyle_AllStatuses(t *testing.T) {
	cases := []Status{
		StatusPending, StatusRunning, StatusPassed,
		StatusFailed, StatusSkipped, StatusErrored,
	}
	for _, s := range cases {
		t.Run(s.String(), func(t *testing.T) {
			style := StatusStyle(s)
			rendered := style.Render(s.String())
			if rendered == "" {
				t.Errorf("StatusStyle(%v).Render returned empty", s)
			}
		})
	}
}

func TestStatusStyle_Unknown(t *testing.T) {
	style := StatusStyle(Status(999))
	if rendered := style.Render("x"); rendered == "" {
		t.Error("unknown status should still render")
	}
}

func TestPaletteRendering(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{"brand", "qube", "qube"},
		{"success", "ok", "ok"},
		{"failure", "ko", "ko"},
		{"warn", "warn", "warn"},
		{"muted", "skip", "skip"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got string
			switch c.name {
			case "brand":
				got = BrandStyle.Render(c.text)
			case "success":
				got = SuccessStyle.Render(c.text)
			case "failure":
				got = FailureStyle.Render(c.text)
			case "warn":
				got = WarnStyle.Render(c.text)
			case "muted":
				got = MutedStyle.Render(c.text)
			}
			if got == "" {
				t.Errorf("%s rendered empty", c.name)
			}
			if !contains(got, c.want) {
				t.Errorf("%s: rendered %q does not contain %q", c.name, got, c.want)
			}
		})
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
