package middleware

import (
	"net/http"

	"github.com/khareutkarshk/dug/internal/httpx"
	"github.com/khareutkarshk/dug/internal/ratelimit"
)

func RateLimit(manager *ratelimit.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			clientIP := httpx.ClientIp(r)

			limiter := manager.Get(clientIP)

			if !limiter.Allow() {
				http.Error(
					w,
					"Rate limit exceeded",
					http.StatusTooManyRequests,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
