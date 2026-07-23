package upstream

import (
	"time"
)

type LeastConnectionsBalancer struct{}

func (LeastConnectionsBalancer) Next(p *Pool) *Backend {

	p.mu.Lock()
	defer p.mu.Unlock()

	var selected *Backend

	for _, backend := range p.backends {

		// skip unhealthy backends
		if !backend.Healthy.Load() {
			continue
		}

		// skip open circuits.
		if backend.CircuitState.Load() == CircuitOpen {
			if time.Now().Unix() < backend.OpenUntil.Load() {
				continue
			}

			if backend.CircuitState.CompareAndSwap(
				CircuitOpen,
				CircuitHalfOpen,
			) {
				backend.EnterHalfOpen()
			}
		}

		// allow only one request to a half-open backend

		if backend.CircuitState.Load() == CircuitHalfOpen {

			if !backend.HalfOpenInFlight.CompareAndSwap(false, true) {
				continue
			}
		}

		if selected == nil {
			selected = backend
			continue
		}

		if backend.ActiveConnections.Load() < selected.ActiveConnections.Load() {
			selected = backend
		}
	}
	return selected
}
