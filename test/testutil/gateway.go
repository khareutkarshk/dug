package testutil

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/khareutkarshk/dug/internal/app"
	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
)

type Gateway struct {
	App *app.App
	URL string
}

func NewGateway(
	t *testing.T,
	opts ...GatewayOptions,
) *Gateway {

	t.Helper()

	options := GatewayOptions{
		RouteTimeout:   5 * time.Second,
		Retries:        2,
		Strategy:       upstream.StrategySmoothWeighted,
		RateLimitRPS:   1000,
		RateLimitBurst: 1000,
	}

	if len(opts) > 0 {
		options = opts[0]
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         0,
			Retries:      options.Retries,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,

			RateLimit: config.RateLimitConfig{
				RPS:   options.RateLimitRPS,
				Burst: options.RateLimitBurst,
			},

			Limits: config.LimitsConfig{
				BodySize: options.BodySize,
			},
			Compression: options.Compression,
			Security: config.SecurityConfig{
				Headers: options.SecurityHeaders,
			},
		},

		Routes: []config.Route{
			{
				Path:    "/",
				Timeout: options.RouteTimeout,

				RequestHeaders: config.HeaderRules{
					Add: map[string]string{
						"X-Gateway": "DUG",
						"X-Version": "v1",
					},
					Remove: []string{
						"X-Internal",
					},
				},

				ResponseHeaders: config.HeaderRules{
					Add: map[string]string{
						"X-Powered-By": "DUG",
					},
					Remove: []string{
						"Server",
					},
				},

				Upstreams: options.Upstreams,

				Strategy: options.Strategy,
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
