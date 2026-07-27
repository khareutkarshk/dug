package upstream

// Pool implements the Smooth Weighted Round Robin algorithm.

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/khareutkarshk/dug/internal/discovery"
)

// Number of consecutive failures before
// a backend is marked unhealthy.
const (
	failureThreshold = 3
)

var (
	CircuitOpenFor = 30 * time.Second
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

	balancer Balancer

	// protects the scheduling algorithm and the backends slice
	mu sync.Mutex
}

func New(provider discovery.Provider) (*Pool, error) {

	pool := &Pool{
		balancer: SmoothWeightedBalancer{},
	}

	upstreams, err := provider.GetUpstreams()
	if err != nil {
		return nil, err
	}

	backends := make([]*Backend, 0, len(upstreams))

	for _, u := range upstreams {

		parsedURL, err := url.Parse(u.URL)
		if err != nil {
			return nil, err
		}

		backend := &Backend{
			URL:    parsedURL,
			Weight: max(u.Weight, 1),
		}

		backend.Healthy.Store(true)
		backend.CircuitState.Store(CircuitClosed)

		backends = append(backends, backend)
		pool.totalWeight += backend.Weight
	}

	pool.backends = backends

	return pool, nil
}

func (p *Pool) Next() *Backend {
	return p.balancer.Next(p)
}

// SetBalancer changes the load-balancing algorithm used by this pool.
func (p *Pool) SetBalancer(b Balancer) {
	p.balancer = b
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
	b.OpenUntil.Store(time.Now().Add(CircuitOpenFor).Unix())
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
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now().Unix()

	for _, backend := range p.backends {

		if !backend.Healthy.Load() {
			continue
		}

		if backend.CircuitState.Load() == CircuitOpen &&
			now < backend.OpenUntil.Load() {
			continue
		}

		return true
	}

	return false
}
