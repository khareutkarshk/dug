package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/khareutkarshk/dug/internal/config"
)

type validateReport struct {
	Valid  bool              `json:"valid"`
	Config string            `json:"config"`
	Server validateServer    `json:"server,omitempty"`
	Routes int               `json:"routes"`
	Issues []ValidationIssue `json:"issues,omitempty"`
}

type validateServer struct {
	Port int  `json:"port"`
	TLS  bool `json:"tls"`
}

// Validate checks a configuration file and prints a human or machine-readable report.
//
// Flags:
//
//	-config  path to YAML config (default configs/edge.yaml)
//	-quiet   print nothing on success; still print errors
//	-json    print a JSON report
func Validate(args []string) error {
	return validate(os.Stdout, os.Stderr, args)
}

func validate(stdout, stderr io.Writer, args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	configPath := fs.String(
		"config",
		"configs/edge.yaml",
		"Path to configuration file",
	)
	quiet := fs.Bool("quiet", false, "Suppress success output")
	jsonOut := fs.Bool("json", false, "Print validation report as JSON")

	if err := fs.Parse(args); err != nil {
		return reported(err)
	}

	absPath, err := filepath.Abs(*configPath)
	if err != nil {
		absPath = *configPath
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		err = fmt.Errorf("load config %s: %w", *configPath, err)
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return reported(err)
	}

	result := checkConfig(cfg)

	rep := validateReport{
		Valid:  result.Ok(),
		Config: absPath,
		Routes: len(cfg.Routes),
		Issues: result.Issues,
	}
	if result.Ok() {
		rep.Server = validateServer{
			Port: cfg.Server.Port,
			TLS:  cfg.Server.TLS.Enabled,
		}
	}

	if *jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			return err
		}
		if !result.Ok() {
			return reported(result.Error())
		}
		return nil
	}

	if !result.Ok() {
		fmt.Fprintln(stderr, "✗ Configuration is invalid")
		fmt.Fprintln(stderr)
		for _, issue := range result.Issues {
			fmt.Fprintf(stderr, "  - %s\n", issue.String())
		}
		return reported(result.Error())
	}

	if *quiet {
		return nil
	}

	fmt.Fprintln(stdout, "✓ Configuration is valid")
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "Config : %s\n", absPath)
	fmt.Fprintf(stdout, "Port   : %d\n", cfg.Server.Port)
	fmt.Fprintf(stdout, "TLS    : %t\n", cfg.Server.TLS.Enabled)
	fmt.Fprintf(stdout, "Routes : %d\n", len(cfg.Routes))
	fmt.Fprintln(stdout)

	for i, route := range cfg.Routes {
		strategy := route.Strategy
		if strategy == "" {
			strategy = "smooth_weighted (default)"
		}
		fmt.Fprintf(
			stdout,
			"  [%d] %s  strategy=%s  upstreams=%d\n",
			i,
			route.Path,
			strategy,
			len(route.Upstreams),
		)
	}

	return nil
}
