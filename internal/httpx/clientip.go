package httpx

import (
	"net"
	"net/http"
	"strings"
)

// ClientIp returns the real client Ip.
// Preference order:
// 1. X-Forwarded-For header
// 2. X-Real-Ip header
// 3. RemoteAddr

func ClientIp(r *http.Request) string {

	// used when DUG is behind a reverse proxy like Nginx or Traefik
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {

		// first value in the X-Forwarded-For header is the original client IP

		parts := strings.Split(xff, ",")

		return strings.TrimSpace(parts[0])
	}

	// used by Ngnix and many load balancers to forward the original client IP
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// fallback to the TCP connection address.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return host
	}

	return r.RemoteAddr
}
