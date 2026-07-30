package cli

import (
	"flag"
	"fmt"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

func Doctor(args []string) error {

	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)

	configPath := fs.String(
		"config",
		"configs/edge.yaml",
		"Configuration file",
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	start := time.Now()

	fmt.Println("Running DUG diagnostics...")
	fmt.Println()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	fmt.Println("Configuration")

	fmt.Println("✓ Configuration loaded")

	if err := validateConfig(cfg); err != nil {
		return err
	}

	fmt.Println("✓ Configuration valid")
	fmt.Println()

	checkPort(cfg.Server.Port)

	fmt.Println()

	checkUpstreams(cfg)

	fmt.Println()

	fmt.Printf("Finished in %s\n", time.Since(start).Round(time.Millisecond))

	return nil
}
