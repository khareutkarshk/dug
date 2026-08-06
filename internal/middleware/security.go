package middleware

import (
	"net/http"
	"strings"

	"github.com/khareutkarshk/dug/internal/config"
)

// SecurityHeaders sets configured response security headers.
// Empty values are skipped so partial configs stay valid.
func SecurityHeaders(cfg config.SecurityHeaders) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setIfPresent(w, "X-Frame-Options", cfg.XFrameOptions)
			setIfPresent(w, "X-Content-Type-Options", cfg.XContentTypeOptions)
			setIfPresent(w, "Strict-Transport-Security", cfg.StrictTransportSecurity)
			setIfPresent(w, "Referrer-Policy", cfg.ReferrerPolicy)
			setIfPresent(w, "Content-Security-Policy", cfg.ContentSecurityPolicy)
			setIfPresent(w, "Permissions-Policy", cfg.PermissionsPolicy)

			next.ServeHTTP(w, r)
		})
	}
}

func setIfPresent(w http.ResponseWriter, name, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	w.Header().Set(name, value)
}

// HasSecurityHeaders reports whether any security header is configured.
func HasSecurityHeaders(cfg config.SecurityHeaders) bool {
	return strings.TrimSpace(cfg.XFrameOptions) != "" ||
		strings.TrimSpace(cfg.XContentTypeOptions) != "" ||
		strings.TrimSpace(cfg.StrictTransportSecurity) != "" ||
		strings.TrimSpace(cfg.ReferrerPolicy) != "" ||
		strings.TrimSpace(cfg.ContentSecurityPolicy) != "" ||
		strings.TrimSpace(cfg.PermissionsPolicy) != ""
}
