package main

import (
	"fmt"
	"log"
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

	case "version", "-v", "--version":
		if err := cli.PrintVersion(args); err != nil {
			log.Fatal(err)
		}

	case "validate":
		if err := cli.Validate(args); err != nil {
			log.Fatal(err)
		}

	case "doctor":
		if err := cli.Doctor(args); err != nil {
			log.Fatal(err)
		}

	case "help", "-h", "--help":
		cli.PrintHelp()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		cli.PrintHelp()
		os.Exit(1)
	}
}
