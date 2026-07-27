package testutil

import "time"

type GatewayOptions struct {
	RouteTimeout time.Duration
	Retries      int
}
