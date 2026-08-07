package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/khareutkarshk/dug/internal/version"
)

// VersionInfo is the structured build metadata exposed by `dug version`.
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	Go        string `json:"go"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// CurrentVersion returns the process build metadata.
func CurrentVersion() VersionInfo {
	v, c, d := version.Info()
	osName, arch := version.Platform()
	return VersionInfo{
		Version:   v,
		Commit:    c,
		BuildDate: d,
		Go:        version.GoVersion(),
		OS:        osName,
		Arch:      arch,
	}
}

// PrintVersion prints version information.
// Supported flags: -short, -json.
func PrintVersion(args []string) error {
	return printVersion(os.Stdout, args)
}

func printVersion(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(w)

	short := fs.Bool("short", false, "Print version number only")
	jsonOut := fs.Bool("json", false, "Print version information as JSON")

	if err := fs.Parse(args); err != nil {
		return err
	}

	info := CurrentVersion()

	switch {
	case *jsonOut:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(info)

	case *short:
		_, err := fmt.Fprintln(w, info.Version)
		return err

	default:
		_, err := fmt.Fprintf(w, `dug %s
  commit:     %s
  build date: %s
  go:         %s
  os/arch:    %s/%s
`, info.Version, info.Commit, info.BuildDate, info.Go, info.OS, info.Arch)
		return err
	}
}
