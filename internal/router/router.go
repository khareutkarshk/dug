package router

import (
	"net/http"
	"time"

	"github.com/khareutkarshk/dug/edge/internal/config"
	"github.com/khareutkarshk/dug/edge/internal/middleware"
	"github.com/khareutkarshk/dug/edge/internal/proxy"
	"github.com/khareutkarshk/dug/edge/internal/upstream"
)

func NewRouter(cfg *config.Config) (http.Handler, error) {

	mux := http.NewServeMux()

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

		handler := middleware.RequireHealthyBackend(pool)(
			middleware.RequestId(
				middleware.Logger(p),
			),
		)

		// register the proxy with the mux
		mux.Handle(route.Path, handler)
	}

	return mux, nil
}
