package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(
		t,
		backend.URL,
		testutil.GatewayOptions{
			RouteTimeout: 1 * time.Second,
		},
	)

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusGatewayTimeout, resp.StatusCode)
}
