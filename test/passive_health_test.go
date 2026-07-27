package test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestPassiveHealthCheck(t *testing.T) {

	var badHits atomic.Int32
	var goodHits atomic.Int32

	// Backend that always fails.
	bad := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		badHits.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	// Healthy backend.
	good := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		goodHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Retries: 0,
			Upstreams: []config.Upstream{
				{URL: bad.URL, Weight: 1},
				{URL: good.URL, Weight: 1},
			},
		},
	)

	// Generate enough requests to trip the circuit breaker.
	for i := 0; i < 10; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		resp.Body.Close()
	}

	require.GreaterOrEqual(t, badHits.Load(), int32(3))
	require.Greater(t, goodHits.Load(), int32(0))
}
