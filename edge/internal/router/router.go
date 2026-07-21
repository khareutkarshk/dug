package router

import (
	"net/http"

	"github.com/khareutkarshk/dug/edge/internal/config"
	"github.com/khareutkarshk/dug/edge/internal/middleware"
	"github.com/khareutkarshk/dug/edge/internal/proxy"
)

func NewRouter(routes []config.Route) (http.Handler, error) {
	mux := http.NewServeMux()

	// loop through the routes and create a proxy for each route
	for _, route := range routes {

		// create a new proxy for the target
		p, err := proxy.New(route.Target)
		if err != nil {
			return nil, err
		}

		handler := middleware.RequestId(
			middleware.Logger(p),
		)

		// register the proxy with the mux
		mux.Handle(route.Path, handler)
	}

	return mux, nil
}
