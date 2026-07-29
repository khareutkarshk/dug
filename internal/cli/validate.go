package cli

import (
	"flag"
	"fmt"

	"github.com/khareutkarshk/dug/internal/config"
)

func Validate(args []string) error {

	fs := flag.NewFlagSet("validate", flag.ContinueOnError)

	configPath := fs.String(
		"config",
		"configs/edge.yaml",
		"Configuration file",
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	fmt.Println("✓ Configuration is valid")
	fmt.Println()

	fmt.Printf("Server Port : %d\n", cfg.Server.Port)
	fmt.Printf("Routes      : %d\n", len(cfg.Routes))

	return nil
}
