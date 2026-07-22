package upstream

import (
	"net/url"
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
}

type Pool struct {
	backends []*Backend
	// weighted round robin schedule
	schedule []*Backend
	current  uint64
}

func New(upstreams []config.Upstream) (*Pool, error) {

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
	}

	// build the weighted round robin schedule

	schedule := make([]*Backend, 0)

	for _, backend := range backends {

		weight := backend.Weight

		// Default to 1 if weight is not set or less than 1
		if weight < 1 {
			weight = 1
		}

		for i := 0; i < weight; i++ {
			schedule = append(schedule, backend)
		}
	}

	return &Pool{
		backends: backends,
		schedule: schedule}, nil
}

func (p *Pool) Next() *Backend {

	if len(p.schedule) == 0 {
		return nil
	}

	start := atomic.AddUint64(&p.current, 1) - 1

	for i := 0; i < len(p.schedule); i++ {

		index := (start + uint64(i)) % uint64(len(p.schedule))

		backend := p.schedule[index]

		if !backend.Healthy.Load() {
			continue
		}

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

		if backend.CircuitState.Load() == CircuitHalfOpen {

			if !backend.HalfOpenInFlight.CompareAndSwap(false, true) {
				continue
			}
		}

		return backend
	}

	return nil
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
