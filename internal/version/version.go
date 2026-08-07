package version

import (
	"runtime"
	"runtime/debug"
)

// These values are intended to be overridden at link time via -ldflags:
//
//	-X github.com/khareutkarshk/dug/internal/version.Version=...
//	-X github.com/khareutkarshk/dug/internal/version.Commit=...
//	-X github.com/khareutkarshk/dug/internal/version.BuildDate=...
//	-X github.com/khareutkarshk/dug/internal/version.Go=...
//
// When unset (typical for `go install`), Info() fills gaps from
// runtime/debug.ReadBuildInfo() so module versions like v0.1.1 appear.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
	Go        = ""
)

// Info returns version metadata, preferring ldflags and falling back to
// Go build info embedded by `go install` / `go build`.
func Info() (version, commit, buildDate string) {
	version, commit, buildDate = Version, Commit, BuildDate

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, commit, buildDate
	}

	if version == "" || version == "dev" {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			version = v
		}
	}

	var vcsRev, vcsTime string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRev = setting.Value
		case "vcs.time":
			vcsTime = setting.Value
		}
	}

	if commit == "" || commit == "unknown" {
		if vcsRev != "" {
			commit = shortCommit(vcsRev)
		}
	}

	if buildDate == "" || buildDate == "unknown" {
		if vcsTime != "" {
			buildDate = vcsTime
		}
	}

	return version, commit, buildDate
}

// GoVersion returns the Go toolchain used to build the binary.
// Falls back to the runtime version when ldflags did not set Go.
func GoVersion() string {
	if Go == "" || Go == "unknown" {
		return runtime.Version()
	}
	return Go
}

// Platform returns the operating system and architecture the binary runs on.
func Platform() (os, arch string) {
	return runtime.GOOS, runtime.GOARCH
}

func shortCommit(rev string) string {
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}
