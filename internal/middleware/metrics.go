package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/khareutkarshk/dug/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		metrics.InFlightRequests.Inc()
		defer metrics.InFlightRequests.Dec()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		start := time.Now()

		next.ServeHTTP(rw, r)

		metrics.RequestDuration.
			WithLabelValues(r.Method, r.URL.Path).
			Observe(time.Since(start).Seconds())

		metrics.RequestsTotal.
			WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(rw.status),
			).
			Inc()
	})
}
