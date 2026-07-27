package testutil

import (
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

type GatewayOptions struct {
	RouteTimeout time.Duration
	Retries      int
	Upstreams    []config.Upstream
	Strategy     string

	RateLimitRPS   float64
	RateLimitBurst int
}
