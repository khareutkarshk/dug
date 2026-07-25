package discovery

import "net/url"

type Endpoint struct {
	URL    *url.URL
	Weight int
}

type Provider interface {
	GetEndpoints() ([]Endpoint, error)
}
