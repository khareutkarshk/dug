package cli

import (
	"errors"
	"net/url"

	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/upstream"
)

func validateConfig(cfg *config.Config) error {

	if cfg.Server.Port <= 0 {
		return errors.New("server.port must be greater than 0")
	}

	for _, route := range cfg.Routes {

		switch route.Strategy {

		case "",
			upstream.StrategyLeastConnections,
			upstream.StrategySmoothWeighted:

		default:
			return errors.New("unknown strategy: " + route.Strategy)
		}

		if len(route.Upstreams) == 0 {
			return errors.New("route has no upstreams")
		}

		for _, u := range route.Upstreams {

			if _, err := url.ParseRequestURI(u.URL); err != nil {
				return err
			}
		}
	}

	return nil
}
