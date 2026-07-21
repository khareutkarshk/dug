package upstream

import (
	"net/url"
	"sync/atomic"
)

type Backned struct {
	URL *url.URL
}

type Pool struct {
	backends []*Backned
	current  uint64
}

func New(upstreams []string) (*Pool, error) {

	backends := make([]*Backned, 0, len(upstreams))

	for _, rawURL := range upstreams {

		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		backends = append(backends, &Backned{URL: u})
	}
	return &Pool{backends: backends}, nil
}

func (p *Pool) Next() *Backned {
	if len(p.backends) == 0 {
		return nil
	}

	index := atomic.AddUint64(&p.current, 1) - 1

	backendIndex := index % uint64(len(p.backends))

	return p.backends[int(backendIndex)]
}
