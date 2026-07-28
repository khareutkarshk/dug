package test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/test/testutil"
	"github.com/stretchr/testify/require"
)

func TestRateLimitBurst(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
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
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

func TestRateLimitRefill(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		RateLimitRPS:   1,
		RateLimitBurst: 1,
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	// Consume the only token.
	resp, err := http.Get(gateway.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
	// Should now be limited.
	resp, err = http.Get(gateway.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
	// Wait for one token to refill.
	time.Sleep(1100 * time.Millisecond)

	resp, err = http.Get(gateway.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
}

func TestRateLimitPerIP(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		RateLimitRPS:   1,
		RateLimitBurst: 1,
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	client := &http.Client{}

	// First request from IP 1 -> allowed
	req, err := http.NewRequest(http.MethodGet, gateway.URL, nil)
	require.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "1.1.1.1")

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
	// Second request from the same IP -> rate limited
	req, err = http.NewRequest(http.MethodGet, gateway.URL, nil)
	require.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "1.1.1.1")

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
	// Different IP should get its own limiter
	req, err = http.NewRequest(http.MethodGet, gateway.URL, nil)
	require.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2.2.2.2")

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer func() {
		_ = resp.Body.Close()
	}()
}

func TestRateLimitConcurrent(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		RateLimitRPS:   1,
		RateLimitBurst: 5,
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	var okCount int64
	var limitedCount int64

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := http.Get(gateway.URL)
			require.NoError(t, err)
			defer func() {
				_ = resp.Body.Close()
			}()
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt64(&okCount, 1)
			case http.StatusTooManyRequests:
				atomic.AddInt64(&limitedCount, 1)
			}
		}()
	}

	wg.Wait()

	require.Equal(t, int64(5), okCount)
	require.Equal(t, int64(15), limitedCount)
	require.Equal(t, int64(20), okCount+limitedCount)
}

func TestRateLimitDisabled(t *testing.T) {
	backend := testutil.NewBackend(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	gateway := testutil.NewGateway(t, testutil.GatewayOptions{
		RateLimitRPS:   0,
		RateLimitBurst: 0,
		Upstreams: []config.Upstream{
			{
				URL:    backend.URL,
				Weight: 1,
			},
		},
	})

	for i := 0; i < 20; i++ {
		resp, err := http.Get(gateway.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer func() {
			_ = resp.Body.Close()
		}()
	}
}
