package test

import (
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestRequestHeaders(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {

		require.Equal(t, "DUG", r.Header.Get("X-Gateway"))
		require.Equal(t, "v1", r.Header.Get("X-Version"))
		require.Empty(t, r.Header.Get("X-Internal"))

		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		testutil.GatewayOptions{
			Upstreams: []config.Upstream{
				{
					URL:    backend.URL,
					Weight: 1,
				},
			},
		},
	)
	req, err := http.NewRequest(http.MethodGet, gateway.URL, nil)
	require.NoError(t, err)

	req.Header.Set("X-Internal", "secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
