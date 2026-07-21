package upstream

import (
	"net/url"
	"sync/atomic"
	"time"
)

// Number of consecutive failures before
// a backend is marked unhealthy.
const failureThreshold = 3

type Backend struct {
	URL *url.URL

	Healthy atomic.Bool

	// Consecutive failed requests.
	Failures atomic.Uint32

	// Time of the most recent failure.
	LastFailure atomic.Int64
}

type Pool struct {
	backends []*Backend
	current  uint64
}

func New(upstreams []string) (*Pool, error) {

	backends := make([]*Backend, 0, len(upstreams))

	for _, upstream := range upstreams {

		u, err := url.Parse(upstream)
		if err != nil {
			return nil, err
		}

		backend := &Backend{
			URL: u,
		}

		backend.Healthy.Store(true)

		backends = append(backends, backend)
	}

	return &Pool{
		backends: backends,
	}, nil
}

func (p *Pool) Next() *Backend {

	if len(p.backends) == 0 {
		return nil
	}

	start := atomic.AddUint64(&p.current, 1) - 1

	for i := 0; i < len(p.backends); i++ {

		index := (start + uint64(i)) % uint64(len(p.backends))

		backend := p.backends[index]

		if backend.Healthy.Load() {
			return backend
		}
	}

	return nil
}

// ReportSuccess is called after a successful request.
func (b *Backend) ReportSuccess() {

	b.Failures.Store(0)

	if !b.Healthy.Load() {
		b.Healthy.Store(true)
	}
}

// ReportFailure is called after a failed request.
func (b *Backend) ReportFailure() {

	failures := b.Failures.Add(1)

	b.LastFailure.Store(time.Now().Unix())

	if failures >= failureThreshold {
		b.Healthy.Store(false)
	}
}

func (p *Pool) HasHealthyBackend() bool {

	for _, backend := range p.backends {

		if backend.Healthy.Load() {
			return true
		}
	}

	return false
}
