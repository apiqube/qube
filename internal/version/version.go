package version

// Values are injected at build time via -ldflags:
//
//	-X github.com/apiqube/qube/internal/version.Version=v1.2.3
//	-X github.com/apiqube/qube/internal/version.Commit=abc1234
//	-X github.com/apiqube/qube/internal/version.Date=2026-04-14
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
