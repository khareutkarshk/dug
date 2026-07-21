package upstream

import (
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL     *url.URL
	Healthy atomic.Bool
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
	return &Pool{backends: backends}, nil
}

// Next returns the next healthy backend using Round Robin.

func (p *Pool) Next() *Backend {
	if len(p.backends) == 0 {
		return nil
	}

	start := atomic.AddUint64(&p.current, 1) - 1

	for i := 0; i < len(p.backends); i++ {
		index := (start + uint64(i)) % uint64(len(p.backends))

		backend := p.backends[int(index)]
		if backend.Healthy.Load() {
			return backend
		}
	}

	return nil
}

// HasHealthyBackend returns true if at least one backend
// is currently healthy.
func (p *Pool) HasHealthyBackend() bool {
	for _, backend := range p.backends {
		if backend.Healthy.Load() {
			return true
		}
	}

	return false
}
