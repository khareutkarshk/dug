package test

import (
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestRateLimitBurst(t *testing.T) {

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		RateLimitRPS:   1,
		RateLimitBurst: 2,
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	for i := 0; i < 2; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}
