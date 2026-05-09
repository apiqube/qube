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
			// Style values are opaque; just verify rendering produces output.
			rendered := style.Render(s.String())
			if rendered == "" {
				t.Errorf("StatusStyle(%v).Render returned empty", s)
			}
		})
	}
}

func TestStatusStyle_Unknown(t *testing.T) {
	// Out-of-range status falls back to muted style.
	style := StatusStyle(Status(999))
	if rendered := style.Render("x"); rendered == "" {
		t.Error("unknown status should still render")
	}
}

func TestPaletteRendering(t *testing.T) {
	// Sanity check: palette styles produce non-empty output.
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
				got = Brand.Render(c.text)
			case "success":
				got = Success.Render(c.text)
			case "failure":
				got = Failure.Render(c.text)
			case "warn":
				got = Warn.Render(c.text)
			case "muted":
				got = Muted.Render(c.text)
			}
			if got == "" {
				t.Errorf("%s rendered empty", c.name)
			}
			// The original text must appear in the output regardless of styling.
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
