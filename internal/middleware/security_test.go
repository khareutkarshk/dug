package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/middleware"
)

func TestSecurityHeadersOnlyConfigured(t *testing.T) {
	h := middleware.SecurityHeaders(config.SecurityHeaders{
		XFrameOptions:       "DENY",
		XContentTypeOptions: "nosniff",
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("X-Frame-Options=%q", got)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q", got)
	}
	if got := rec.Header().Get("Referrer-Policy"); got != "" {
		t.Fatalf("unexpected Referrer-Policy=%q", got)
	}
}
