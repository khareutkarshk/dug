package cli

import "fmt"

func PrintHelp() {
	fmt.Println("DUG - Distributed Unified Gateway")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("  dug <command> [options]")
	fmt.Println()

	fmt.Println("Commands:")
	fmt.Println("  run        Start the gateway")
	fmt.Println("  validate   Validate configuration")
	fmt.Println("  version    Print version information")
	fmt.Println("  help       Show this help message")
	fmt.Println()

	fmt.Println("Examples:")
	fmt.Println("  dug run -config configs/edge.yaml")
	fmt.Println("  dug validate -config configs/edge.yaml")
	fmt.Println("  dug version")
}
