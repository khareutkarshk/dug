package discovery

import (
	"net/url"

	"github.com/khareutkarshk/dug/internal/config"
)

type Static struct {
	Upstreams []config.Upstream
}

func (s Static) GetEndpoints() ([]Endpoint, error) {

	endpoints := make([]Endpoint, 0, len(s.Upstreams))

	for _, u := range s.Upstreams {

		parsed, err := url.Parse(u.URL)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, Endpoint{
			URL:    parsed,
			Weight: u.Weight,
		})
	}

	return endpoints, nil
}
