//go:build benchmark

package benchmarks

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/khareutkarshk/dug/internal/app"
	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/logger"
	"github.com/khareutkarshk/dug/internal/upstream"
)

// Harness runs an in-process DUG gateway in front of httptest upstreams.
type Harness struct {
	GatewayURL string
	Strategy   string

	app      *app.App
	backends []*httptest.Server
}

// StartHarness boots equal-weight httptest backends and a DUG gateway using the
// given load-balancing strategy.
func StartHarness(strategy string, backendCount int) (*Harness, error) {
	if backendCount <= 0 {
		backendCount = 3
	}

	// Quiet request logs so they do not dominate measured latency.
	logger.Log = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	backends := make([]*httptest.Server, 0, backendCount)
	upstreams := make([]config.Upstream, 0, backendCount)

	for i := 0; i < backendCount; i++ {
		backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		backends = append(backends, backend)
		upstreams = append(upstreams, config.Upstream{
			URL:    backend.URL,
			Weight: 1,
		})
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			Retries:      0,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
			RateLimit: config.RateLimitConfig{
				// Effectively disable rate limiting for load measurement.
				RPS:   1_000_000,
				Burst: 1_000_000,
			},
		},
		Routes: []config.Route{{
			Path:      "/",
			Timeout:   5 * time.Second,
			Strategy:  strategy,
			Upstreams: upstreams,
		}},
	}

	edge, err := app.New(cfg)
	if err != nil {
		closeBackends(backends)
		return nil, fmt.Errorf("create gateway: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		err := edge.Server.Start()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-edge.Server.Ready():
	case err := <-errCh:
		closeBackends(backends)
		return nil, fmt.Errorf("start gateway: %w", err)
	case <-time.After(5 * time.Second):
		closeBackends(backends)
		return nil, fmt.Errorf("gateway did not become ready")
	}

	return &Harness{
		GatewayURL: "http://" + edge.Server.Addr(),
		Strategy:   strategy,
		app:        edge,
		backends:   backends,
	}, nil
}

// Close shuts down the gateway and upstream servers.
func (h *Harness) Close() error {
	if h == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var firstErr error
	if h.app != nil {
		if err := h.app.Server.Shutdown(ctx); err != nil {
			firstErr = err
		}
	}
	closeBackends(h.backends)
	return firstErr
}

func closeBackends(backends []*httptest.Server) {
	for _, b := range backends {
		if b != nil {
			b.Close()
		}
	}
}

// SupportedStrategies are the load-balancing strategies exercised by the suite.
func SupportedStrategies() []string {
	return []string{
		upstream.StrategySmoothWeighted,
		upstream.StrategyLeastConnections,
	}
}
