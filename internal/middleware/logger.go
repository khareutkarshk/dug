package middleware

import (
	"net/http"
	"time"

	"github.com/khareutkarshk/dug/internal/httpx"
	"github.com/khareutkarshk/dug/internal/logger"
)

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, _ := r.Context().Value(RequestIDKey).(string)

		start := time.Now()

		rw := httpx.NewResponseWriter(w)

		next.ServeHTTP(w, r)

		logger.Log.Info(
			"http request",
			"request_id", id,
			"method", r.Method,
			"status", rw.StatusCode,
			"bytes_written", rw.BytesWritten,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
