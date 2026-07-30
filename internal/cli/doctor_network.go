package cli

import (
	"fmt"
	"net"
)

func checkPort(port int) {

	fmt.Println("Gateway")

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {

		fmt.Printf("✗ Port %d already in use\n", port)
		return
	}

	_ = ln.Close()

	fmt.Printf("✓ Port %d available\n", port)
}
