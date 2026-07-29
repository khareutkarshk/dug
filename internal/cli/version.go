package cli

import (
	"fmt"

	"github.com/khareutkarshk/dug/internal/version"
)

func PrintVersion() {
	fmt.Println("DUG")
	fmt.Printf("Version: %s\n", version.Version)
}
