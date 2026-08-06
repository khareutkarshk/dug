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
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Go      string `json:"go"`
}

// CurrentVersion returns the process build metadata.
func CurrentVersion() VersionInfo {
	v, c, d := version.Info()
	return VersionInfo{
		Version: v,
		Commit:  c,
		Date:    d,
		Go:      version.GoVersion(),
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
  commit: %s
  built:  %s
  go:     %s
`, info.Version, info.Commit, info.Date, info.Go)
		return err
	}
}
