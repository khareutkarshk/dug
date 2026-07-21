package router

import (
	"net/http"

	"github.com/khareutkarshk/dug/edge/internal/config"
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

		// register the proxy with the mux
		mux.Handle(route.Path, p)
	}

	return mux, nil
}
