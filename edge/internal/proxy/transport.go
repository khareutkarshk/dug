package proxy

import (
	"net/http"

	"github.com/khareutkarshk/dug/edge/internal/upstream"
)

// RetryTransport retries failed requests against another backend.
type RetryTransport struct {
	Base    http.RoundTripper
	Pool    *upstream.Pool
	Retries int
}

// Methods that are safe to retry.
var retryableMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	// Non-idempotent requests are forwarded only once.
	if !retryableMethods[req.Method] {

		backend := t.Pool.Next()
		if backend == nil {
			return nil, http.ErrServerClosed
		}

		r := req.Clone(req.Context())
		r.URL.Scheme = backend.URL.Scheme
		r.URL.Host = backend.URL.Host
		r.Host = backend.URL.Host

		return base.RoundTrip(r)
	}

	var lastErr error

	for attempt := 0; attempt <= t.Retries; attempt++ {

		backend := t.Pool.Next()
		if backend == nil {
			break
		}

		r := req.Clone(req.Context())
		r.URL.Scheme = backend.URL.Scheme
		r.URL.Host = backend.URL.Host
		r.Host = backend.URL.Host

		resp, err := base.RoundTrip(r)

		// Retry on network failures.
		if err != nil {
			backend.Healthy.Store(false)
			lastErr = err
			continue
		}

		// Retry on temporary upstream failures.
		switch resp.StatusCode {
		case http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:

			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return &http.Response{
		StatusCode: http.StatusBadGateway,
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}
