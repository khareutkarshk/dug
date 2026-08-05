package cli

import "fmt"

func PrintHelp() {
	fmt.Println(`DUG - Distributed Unified Gateway

Usage:
  dug <command> [options]

Commands:
  run         Start the gateway
  validate    Validate a configuration file
  doctor      Run local diagnostics against a config
  version     Print version information
  help        Show this help message

Global flags:
  -h, --help       Show help
  -v, --version    Print version information

Run options:
  -config string   Path to configuration file (default "configs/edge.yaml")

Validate / doctor options:
  -config string   Path to configuration file (default "configs/edge.yaml")

Version options:
  -short           Print version number only
  -json            Print version information as JSON

Examples:
  dug run -config configs/edge.yaml
  dug validate -config configs/edge.yaml
  dug doctor -config configs/edge.yaml
  dug version
  dug version -short
  dug version -json

Install:
  go install github.com/khareutkarshk/dug/cmd/dug@latest`)
}
