package test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestLeastConnections(t *testing.T) {

	var slowHits atomic.Int32
	var fastHits atomic.Int32

	slow := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		slowHits.Add(1)

		time.Sleep(500 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
	})

	fast := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

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

	for i := 0; i < 20; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get(gateway.URL)
			require.NoError(t, err)
			resp.Body.Close()
		}()
	}

	wg.Wait()

	require.Greater(t, fastHits.Load(), slowHits.Load())
}
