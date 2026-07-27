package test

import (
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestRequestHeaders(t *testing.T) {
	var gatewayHeader string
	var versionHeader string
	var internalHeader string

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		gatewayHeader = r.Header.Get("X-Gateway")
		versionHeader = r.Header.Get("X-Version")
		internalHeader = r.Header.Get("X-Internal")

		w.WriteHeader(http.StatusOK)
	})
	defer backend.Close()

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	req, err := http.NewRequest(http.MethodGet, gateway.URL, nil)
	require.NoError(t, err)

	req.Header.Set("X-Internal", "secret")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Equal(t, "DUG", gatewayHeader)
	require.Equal(t, "v1", versionHeader)
	require.Empty(t, internalHeader)
}
