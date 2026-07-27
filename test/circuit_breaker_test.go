package test

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreakerRecovery(t *testing.T) {

	old := upstream.CircuitOpenFor
	upstream.CircuitOpenFor = 200 * time.Millisecond
	defer func() {
		upstream.CircuitOpenFor = old
	}()

	var fail atomic.Bool
	fail.Store(true)

	var hits atomic.Int32

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		hits.Add(1)

		if fail.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Retries: 0,
			Upstreams: []config.Upstream{
				{
					URL:    backend.URL,
					Weight: 1,
				},
			},
		},
	)

	// Trip the circuit.
	for i := 0; i < 3; i++ {
		t.Log("making recovery request")
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		resp.Body.Close()
	}

	// Wait for the circuit to become half-open.
	time.Sleep(250 * time.Millisecond)

	// Backend becomes healthy.
	fail.Store(false)

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)

	t.Log("status:", resp.StatusCode)

	resp.Body.Close()

	require.Equal(t, int32(4), hits.Load())
}
