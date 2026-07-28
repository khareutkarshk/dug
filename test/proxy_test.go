package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestReverseProxy(t *testing.T) {

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
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
	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "hello", string(body))
}
