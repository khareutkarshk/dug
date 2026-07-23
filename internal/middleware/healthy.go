package middleware

import (
	"net/http"

	"github.com/khareutkarshk/dug/internal/upstream"
)

// HealthyBackend ensures that at least one backend is available.
func RequireHealthyBackend(pool *upstream.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !pool.HasHealthyBackend() {
				http.Error(
					w,
					"No healthy upstreams",
					http.StatusServiceUnavailable,
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
