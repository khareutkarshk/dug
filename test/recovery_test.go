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

func TestBackendRecovery(t *testing.T) {

	var healthy atomic.Int32
	var recovered atomic.Int32

	var backendHealthy atomic.Bool
	backendHealthy.Store(false)

	backend1 := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		healthy.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	backend2 := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {

			if backendHealthy.Load() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}

			return
		}

		recovered.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Upstreams: []config.Upstream{
				{
					URL:    backend1.URL,
					Weight: 1,
				},
				{
					URL:    backend2.URL,
					Weight: 1,
				},
			},
		},
	)

	// Wait for first health check.
	time.Sleep(6 * time.Second)

	// Initially all traffic should go to backend1.
	for i := 0; i < 10; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	require.Equal(t, int32(10), healthy.Load())
	require.Equal(t, int32(0), recovered.Load())

	// Recover backend2.
	backendHealthy.Store(true)

	// Wait for next health check cycle.
	time.Sleep(6 * time.Second)

	healthy.Store(0)
	recovered.Store(0)

	// Now both backends should receive traffic.
	for i := 0; i < 20; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	require.Greater(t, healthy.Load(), int32(0))
	require.Greater(t, recovered.Load(), int32(0))
}
