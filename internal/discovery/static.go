package discovery

import "github.com/khareutkarshk/dug/internal/config"

type Static struct {
	Upstreams []config.Upstream
}

func (s Static) GetUpstreams() ([]config.Upstream, error) {
	return s.Upstreams, nil
}
