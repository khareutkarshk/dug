package upstream

// Pool implements the Smooth Weighted Round Robin algorithm.

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/khareutkarshk/dug/internal/config"
)

// Number of consecutive failures before
// a backend is marked unhealthy.
const (
	failureThreshold = 3
	circuitOpenFor   = 30 * time.Second
)

const (
	CircuitClosed uint32 = iota
	CircuitOpen
	CircuitHalfOpen
)

type Backend struct {
	URL     *url.URL
	Weight  int
	Healthy atomic.Bool

	// Consecutive failed requests.
	Failures atomic.Uint32

	// Time of the most recent failure.
	LastFailure atomic.Int64

	CircuitState     atomic.Uint32
	OpenUntil        atomic.Int64
	HalfOpenInFlight atomic.Bool

	// used only by the load balancing algorithm
	// access must be protected by Pool.mu
	CurrentWeight int

	// Number of active in-flight requests to this backend.
	// used by the Least Connections load balancing algorithm
	ActiveConnections atomic.Int64
}

type Pool struct {
	backends    []*Backend
	totalWeight int

	// protects the scheduling algorithm and the backends slice
	mu sync.Mutex
}

func New(upstreams []config.Upstream) (*Pool, error) {

	pool := &Pool{}

	backends := make([]*Backend, 0, len(upstreams))

	for _, upstream := range upstreams {

		u, err := url.Parse(upstream.URL)
		if err != nil {
			return nil, err
		}

		backend := &Backend{
			URL:    u,
			Weight: max(upstream.Weight, 1), // Ensure weight is at least 1
		}

		backend.CircuitState.Store(CircuitClosed)

		backend.Healthy.Store(true)

		backends = append(backends, backend)

		pool.totalWeight += backend.Weight
	}

	// build the weighted round robin schedule

	pool.backends = backends

	return pool, nil
}

func (p *Pool) Next() *Backend {

	p.mu.Lock()
	defer p.mu.Unlock()

	var selected *Backend

	for _, backend := range p.backends {

		// skip unhealthy backends
		if !backend.Healthy.Load() {
			continue
		}

		// Handle open circuts

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

// ReportSuccess is called after a successful request.
func (b *Backend) ReportSuccess() {

	if b.CircuitState.Load() == CircuitHalfOpen {
		b.CloseCircuit()
		return
	}
	b.Failures.Store(0)
}

// ReportFailure is called after a failed request.
func (b *Backend) ReportFailure() {

	if b.CircuitState.Load() == CircuitHalfOpen {
		b.OpenCircuit()
		b.Failures.Store(0)
		return
	}

	failures := b.Failures.Add(1)

	if failures >= failureThreshold {
		b.OpenCircuit()
	} else {
		b.LastFailure.Store(time.Now().Unix())
	}
}

func (b *Backend) OpenCircuit() {
	b.Healthy.Store(false)
	b.CircuitState.Store(CircuitOpen)
	b.OpenUntil.Store(time.Now().Add(circuitOpenFor).Unix())
	b.HalfOpenInFlight.Store(false)
}

func (b *Backend) CloseCircuit() {
	b.CircuitState.Store(CircuitClosed)
	b.HalfOpenInFlight.Store(false)
	b.Healthy.Store(true)
	b.Failures.Store(0)
}

func (b *Backend) EnterHalfOpen() {
	b.CircuitState.Store(CircuitHalfOpen)
	b.HalfOpenInFlight.Store(false)
}

func (p *Pool) HasHealthyBackend() bool {

	for _, backend := range p.backends {

		if !backend.Healthy.Load() {
			continue
		}

		if backend.CircuitState.Load() == CircuitOpen {
			continue
		}

		return true
	}

	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
