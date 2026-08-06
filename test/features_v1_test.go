package test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestBodySizeLimitIntegration(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		Upstreams: []config.Upstream{{URL: backend.URL, Weight: 1}},
		BodySize:  16,
	})

	req, err := http.NewRequest(http.MethodPost, gateway.URL+"/", strings.NewReader(strings.Repeat("a", 64)))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestCompressionIntegration(t *testing.T) {
	payload := strings.Repeat("benchmark-data-", 200)

	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(payload))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		Upstreams: []config.Upstream{{URL: backend.URL, Weight: 1}},
		Compression: config.CompressionConfig{
			Enabled: true,
			MinSize: 100,
		},
	})

	req, err := http.NewRequest(http.MethodGet, gateway.URL+"/", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{
		Transport: &http.Transport{DisableCompression: true},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gr, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	defer func() { _ = gr.Close() }()

	body, err := io.ReadAll(gr)
	require.NoError(t, err)
	require.Equal(t, payload, string(body))
}

func TestSecurityHeadersIntegration(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		Upstreams: []config.Upstream{{URL: backend.URL, Weight: 1}},
		SecurityHeaders: config.SecurityHeaders{
			XFrameOptions:       "DENY",
			XContentTypeOptions: "nosniff",
			ReferrerPolicy:      "no-referrer",
		},
	})

	resp, err := http.Get(gateway.URL + "/")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	require.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	require.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
	require.Empty(t, resp.Header.Get("Content-Security-Policy"))
}

func TestCompressionDisabledByDefault(t *testing.T) {
	payload := bytes.Repeat([]byte("z"), 2048)
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		Upstreams: []config.Upstream{{URL: backend.URL, Weight: 1}},
	})

	req, err := http.NewRequest(http.MethodGet, gateway.URL+"/", nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{
		Transport: &http.Transport{DisableCompression: true},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	require.Empty(t, resp.Header.Get("Content-Encoding"))
}
