package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/discovery"
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

		provider := discovery.Static{
			Upstreams: route.Upstreams,
		}

		pool, err := upstream.New(provider)
		if err != nil {
			return nil, err
		}

		switch route.Strategy {
		case "", upstream.StrategySmoothWeighted:
			pool.SetBalancer(upstream.SmoothWeightedBalancer{})

		case upstream.StrategyLeastConnections:
			pool.SetBalancer(upstream.LeastConnectionsBalancer{})

		default:
			return nil, fmt.Errorf("unknown load balancing strategy: %s", route.Strategy)
		}

		// start the health check for the upstreams in background
		pool.StartHealthCheck(5 * time.Second)

		p := proxy.New(
			pool,
			cfg.Server.Retries,
			route.Timeout,
			route.RequestHeaders,
			route.ResponseHeaders,
		)

		handler := http.Handler(p)

		handler = middleware.Metrics(handler)
		handler = middleware.Logger(handler)
		handler = middleware.RequestId(handler)

		if cfg.Server.RateLimit.RPS > 0 && cfg.Server.RateLimit.Burst > 0 {
			manager := ratelimit.NewManager(
				cfg.Server.RateLimit.RPS,
				cfg.Server.RateLimit.Burst,
			)
			handler = middleware.RateLimit(manager)(handler)
		}

		handler = middleware.CORS(route.CORS)(handler)

		// register the proxy with the mux
		mux.Handle(route.Path, handler)
	}

	return mux, nil
}
