package test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestLeastConnections(t *testing.T) {
	var slowHits atomic.Int32
	var fastHits atomic.Int32

	slow := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		slowHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	fast := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		fastHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Strategy: upstream.StrategyLeastConnections,
			Upstreams: []config.Upstream{
				{
					URL:    slow.URL,
					Weight: 1,
				},
				{
					URL:    fast.URL,
					Weight: 1,
				},
			},
		},
	)

	var wg sync.WaitGroup

	const requests = 20

	for i := 0; i < requests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get(gateway.URL)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	wg.Wait()

	require.EqualValues(
		t,
		requests,
		slowHits.Load()+fastHits.Load(),
	)

	require.Greater(t, slowHits.Load(), int32(0))
	require.Greater(t, fastHits.Load(), int32(0))
}
