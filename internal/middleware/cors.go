package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/khareutkarshk/dug/internal/config"
)

func CORS(cfg config.CORSConfig) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")

			switch {
			case len(cfg.AllowOrigins) == 0:
				// Do nothing

			case len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*":
				w.Header().Set(
					"Access-Control-Allow-Origin",
					"*",
				)

			default:
				for _, allowed := range cfg.AllowOrigins {
					if origin == allowed {
						w.Header().Set(
							"Access-Control-Allow-Origin",
							origin,
						)
						break
					}
				}
			}

			w.Header().Add("Vary", "Origin")

			if len(cfg.AllowMethods) > 0 {
				w.Header().Set(
					"Access-Control-Allow-Methods",
					strings.Join(cfg.AllowMethods, ","),
				)
			}

			if len(cfg.AllowHeaders) > 0 {
				w.Header().Set(
					"Access-Control-Allow-Headers",
					strings.Join(cfg.AllowHeaders, ","),
				)
			}

			if len(cfg.ExposeHeaders) > 0 {
				w.Header().Set(
					"Access-Control-Expose-Headers",
					strings.Join(cfg.ExposeHeaders, ","),
				)
			}

			if cfg.AllowCredentials {
				w.Header().Set(
					"Access-Control-Allow-Credentials",
					"true",
				)
			}

			if cfg.MaxAge > 0 {
				w.Header().Set(
					"Access-Control-Max-Age",
					strconv.Itoa(cfg.MaxAge),
				)
			}

			// Handle browser preflight
			if r.Method == http.MethodOptions &&
				r.Header.Get("Access-Control-Request-Method") != "" {

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
