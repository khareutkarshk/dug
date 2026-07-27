package testutil

import (
	"net"
	"time"
)

func waitForServer(addr string) {

	for i := 0; i < 50; i++ {

		if addr != "" {

			conn, err := net.DialTimeout(
				"tcp",
				addr,
				100*time.Millisecond,
			)

			if err == nil {
				conn.Close()
				return
			}
		}

		time.Sleep(20 * time.Millisecond)
	}
}
