package router

import (
	"net/http"

	"github.com/khareutkarshk/dug/edge/internal/config"
	"github.com/khareutkarshk/dug/edge/internal/middleware"
	"github.com/khareutkarshk/dug/edge/internal/proxy"
	"github.com/khareutkarshk/dug/edge/internal/upstream"
)

func NewRouter(routes []config.Route) (http.Handler, error) {
	mux := http.NewServeMux()

	// loop through the routes and create a proxy for each route
	for _, route := range routes {

		// create a new proxy for the target
		pool, err := upstream.New(route.Upstreams)
		if err != nil {
			return nil, err
		}

		p := proxy.New(pool)

		handler := middleware.RequestId(
			middleware.Logger(p),
		)

		// register the proxy with the mux
		mux.Handle(route.Path, handler)
	}

	return mux, nil
}
