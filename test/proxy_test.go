package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestReverseProxy(t *testing.T) {

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	gateway := testutil.NewGateway(t, backend.URL)

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "hello", string(body))
}
