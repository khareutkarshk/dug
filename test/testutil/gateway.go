package testutil

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/app"
	"github.com/khareutkarshk/dug/internal/config"
)

type Gateway struct {
	App *app.App
	URL string
}

func NewGateway(t *testing.T, upstream string) *Gateway {
	t.Helper()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,

			RateLimit: config.RateLimitConfig{
				RPS:   1000,
				Burst: 1000,
			},
		},
		Routes: []config.Route{
			{
				Path: "/",

				RequestHeaders: config.HeaderRules{
					Add: map[string]string{
						"X-Gateway": "DUG",
						"X-Version": "v1",
					},
					Remove: []string{"X-Internal"},
				},
				Upstreams: []config.Upstream{
					{
						URL:    upstream,
						Weight: 1,
					},
				},
			},
		},
	}

	edge, err := app.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := edge.Server.Start(); err != nil &&
			err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	// Wait until the server is actually listening.
	<-edge.Server.Ready()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_ = edge.Server.Shutdown(ctx)
	})

	return &Gateway{
		App: edge,
		URL: "http://" + edge.Server.Addr(),
	}
}
