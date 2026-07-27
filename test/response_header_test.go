package test

import (
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestResponseHeaders(t *testing.T) {

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "nginx")
		w.Header().Set("X-Test", "backend")
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(t, backend.URL)

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, "DUG", resp.Header.Get("X-Powered-By"))
	require.Empty(t, resp.Header.Get("Server"))

	require.Equal(t, "backend", resp.Header.Get("X-Test"))
}
