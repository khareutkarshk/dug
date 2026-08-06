package cli

import (
	"fmt"
	"os"
)

func PrintHelp() {
	_, _ = fmt.Fprintln(os.Stdout, `DUG - Distributed Unified Gateway

Usage:
  dug <command> [options]

Commands:
  init        Bootstrap a new DUG project
  run         Start the gateway
  validate    Validate a configuration file
  doctor      Run local diagnostics against a config
  version     Print version information
  help        Show this help message

Global flags:
  -h, --help       Show help
  -v, --version    Print version information

Init options:
  dug init <directory>
  --force          Overwrite existing files

Run options:
  -config string   Path to configuration file (default "configs/edge.yaml")

Validate options:
  -config string   Path to configuration file (default "configs/edge.yaml")
  -quiet           Suppress success output (still prints errors)
  -json            Print validation report as JSON

Doctor options:
  -config string   Path to configuration file (default "configs/edge.yaml")
  -timeout duration
                   Per-upstream probe timeout (default 3s)
  -json            Print diagnostics as JSON

Version options:
  -short           Print version number only
  -json            Print version information as JSON

Examples:
  dug init my-gateway
  dug run -config configs/edge.yaml
  dug validate -config configs/edge.yaml
  dug validate -config configs/edge.yaml -json
  dug doctor -config configs/edge.yaml
  dug doctor -config configs/edge.yaml -timeout 2s -json
  dug version
  dug version -short
  dug version -json

Install:
  go install github.com/khareutkarshk/dug/cmd/dug@latest`)
}
