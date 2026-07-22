package proxy

import (
	"errors"
	"net/http"
	"time"

	"github.com/khareutkarshk/dug/edge/internal/upstream"
)

var (
	ErrNoHealthyBackend = errors.New("no healthy backend available")
)

var retryableMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodOptions: {},
}

// Base delay for retries.
const initialBackoff = 100 * time.Millisecond

type RetryTransport struct {
	Base    http.RoundTripper
	Pool    *upstream.Pool
	Retries int
}

// RoundTrip is called by ReverseProxy for every outgoing request.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	maxAttempts := 1

	// Only retry idempotent requests.
	if _, ok := retryableMethods[req.Method]; ok {
		maxAttempts += t.Retries
	}

	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {

		resp, err := t.send(base, req)

		if err == nil {

			switch resp.StatusCode {
			case http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout:

				resp.Body.Close()

			default:
				return resp, nil
			}

		} else {
			lastErr = err
		}

		// Don't sleep after the final attempt.
		if attempt < maxAttempts-1 {

			// Exponential backoff:
			// 100ms → 200ms → 400ms → 800ms ...
			backoff := initialBackoff * time.Duration(1<<attempt)

			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(backoff):
			}
		}
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

// send forwards a request to one backend.
func (t *RetryTransport) send(
	base http.RoundTripper,
	req *http.Request,
) (*http.Response, error) {

	backend := t.Pool.Next()
	if backend == nil {
		return nil, ErrNoHealthyBackend
	}

	r := req.Clone(req.Context())

	r.URL.Scheme = backend.URL.Scheme
	r.URL.Host = backend.URL.Host
	r.Host = backend.URL.Host

	resp, err := base.RoundTrip(r)

	if err != nil {
		backend.ReportFailure()
		return nil, err
	}

	// Treat temporary server errors as failures.
	switch resp.StatusCode {
	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		backend.ReportFailure()

	default:
		backend.ReportSuccess()
	}

	return resp, nil
}
