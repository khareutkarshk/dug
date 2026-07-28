package test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	var attempts atomic.Int32

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		n := attempts.Add(1)

		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Retries: 2,
			Upstreams: []config.Upstream{
				{
					URL:    backend.URL,
					Weight: 1,
				},
			},
		},
	)
	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, int32(3), attempts.Load())
}
