package test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestSmoothWeightedRoundRobin(t *testing.T) {

	var backend1Hits atomic.Int32
	var backend2Hits atomic.Int32

	backend1 := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		backend1Hits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	backend2 := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		backend2Hits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Upstreams: []config.Upstream{
				{
					URL:    backend1.URL,
					Weight: 3,
				},
				{
					URL:    backend2.URL,
					Weight: 1,
				},
			},
		},
	)

	for i := 0; i < 40; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		resp.Body.Close()
	}

	require.Equal(t, int32(30), backend1Hits.Load())
	require.Equal(t, int32(10), backend2Hits.Load())
}
