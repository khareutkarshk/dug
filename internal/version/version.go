package version

import "runtime"

// These values are intended to be overridden at link time via -ldflags:
//
//	-X github.com/khareutkarshk/dug/internal/version.Version=...
//	-X github.com/khareutkarshk/dug/internal/version.Commit=...
//	-X github.com/khareutkarshk/dug/internal/version.Date=...
//	-X github.com/khareutkarshk/dug/internal/version.Go=...
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
	Go      = ""
)

// GoVersion returns the Go toolchain used to build the binary.
// Falls back to the runtime version when ldflags did not set Go.
func GoVersion() string {
	if Go == "" || Go == "unknown" {
		return runtime.Version()
	}
	return Go
}
