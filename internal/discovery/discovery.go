package discovery

import "github.com/khareutkarshk/dug/internal/config"

type Provider interface {
	GetUpstreams() ([]config.Upstream, error)
}
