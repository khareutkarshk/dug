package main

import (
	"fmt"
	"log"
	"os"

	"github.com/khareutkarshk/dug/internal/cli"
)

func main() {

	if len(os.Args) < 2 {
		run(os.Args[1:])
		return
	}

	switch os.Args[1] {

	case "run":
		run(os.Args[2:])

	case "version":
		cli.PrintVersion()

	case "validate":
		if err := cli.Validate(os.Args[2:]); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
