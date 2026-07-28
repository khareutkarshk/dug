package test

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestHealthyBackendOnly(t *testing.T) {

	var healthyHits atomic.Int32
	var unhealthyHits atomic.Int32

	healthy := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		healthyHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	unhealthy := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		unhealthyHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Upstreams: []config.Upstream{
				{
					URL:    healthy.URL,
					Weight: 1,
				},
				{
					URL:    unhealthy.URL,
					Weight: 1,
				},
			},
		},
	)

	// Wait for one health-check cycle.
	time.Sleep(6 * time.Second)

	for i := 0; i < 20; i++ {

		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)

		defer func() {
			_ = resp.Body.Close()
		}()
	}

	require.Equal(t, int32(20), healthyHits.Load())
	require.Equal(t, int32(0), unhealthyHits.Load())
}
