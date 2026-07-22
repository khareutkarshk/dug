package router

import (
	"net/http"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/middleware"
	"github.com/khareutkarshk/dug/internal/proxy"
	"github.com/khareutkarshk/dug/internal/ratelimit"
	"github.com/khareutkarshk/dug/internal/upstream"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(cfg *config.Config) (http.Handler, error) {

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	// loop through the routes and create a proxy for each route
	for _, route := range cfg.Routes {

		// create a new proxy for the target
		pool, err := upstream.New(route.Upstreams)
		if err != nil {
			return nil, err
		}
		// start the health check for the upstreams in background
		pool.StartHealthCheck(5 * time.Second)

		p := proxy.New(pool, cfg.Server.Retries)

		manager := ratelimit.NewManager(
			cfg.Server.RateLimit.RPS,
			cfg.Server.RateLimit.Burst,
		)

		handler := middleware.RequireHealthyBackend(pool)(
			middleware.RateLimit(manager)(
				middleware.RequestId(
					middleware.Logger(
						middleware.Metrics(p),
					),
				),
			),
		)

		// register the proxy with the mux
		mux.Handle(route.Path, handler)
	}

	return mux, nil
}
