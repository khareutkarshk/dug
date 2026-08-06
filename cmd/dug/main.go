package main

import (
	"fmt"
	"os"

	"github.com/khareutkarshk/dug/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		run(nil)
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "run":
		run(args)

	case "init":
		exit(cli.Init(args))

	case "version", "-v", "--version":
		exit(cli.PrintVersion(args))

	case "validate":
		exit(cli.Validate(args))

	case "doctor":
		exit(cli.Doctor(args))

	case "help", "-h", "--help":
		cli.PrintHelp()

	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		cli.PrintHelp()
		os.Exit(1)
	}
}

func exit(err error) {
	if err == nil {
		return
	}
	if !cli.IsReported(err) {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
