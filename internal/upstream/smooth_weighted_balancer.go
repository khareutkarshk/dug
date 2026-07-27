// Smooth Weighted Round Robin load balancer.
//
// Reference:
// https://github.com/phusion/nginx/blob/master/src/http/ngx_http_upstream_round_robin.c

package upstream

import (
	"time"
)

type SmoothWeightedBalancer struct{}

func (SmoothWeightedBalancer) Next(p *Pool) *Backend {

	p.mu.Lock()
	defer p.mu.Unlock()

	var selected *Backend

	for _, backend := range p.backends {

		state := backend.CircuitState.Load()

		// Handle open circuits first.
		if state == CircuitOpen {

			if time.Now().Unix() < backend.OpenUntil.Load() {
				continue
			}

			if backend.CircuitState.CompareAndSwap(
				CircuitOpen,
				CircuitHalfOpen,
			) {
				backend.EnterHalfOpen()
			}

			state = CircuitHalfOpen
		}

		// Allow only one request to a half-open backend.
		if state == CircuitHalfOpen {

			if !backend.HalfOpenInFlight.CompareAndSwap(false, true) {
				continue
			}
		}

		// Skip unhealthy backends only if they are not part of the
		// circuit breaker flow.
		if !backend.Healthy.Load() && state == CircuitClosed {
			continue
		}

		// smooth weighted round robin algorithm
		// Increase the current weight by configured weight
		backend.CurrentWeight += backend.Weight

		// Select the backend with the highest current weight
		if selected == nil || backend.CurrentWeight > selected.CurrentWeight {
			selected = backend
		}
	}

	if selected == nil {
		return nil
	}

	// reduce the current weight of the selected backend by the total weight
	selected.CurrentWeight -= p.totalWeight

	return selected
}
