package discovery

import (
	"fmt"
	"net"

	"github.com/khareutkarshk/dug/internal/config"
)

type DNS struct {
	Scheme string
	Host   string
	Port   int
}

func (d DNS) GetUpstreams() ([]config.Upstream, error) {

	ips, err := net.LookupIP(d.Host)
	if err != nil {
		return nil, err
	}

	upstreams := make([]config.Upstream, 0, len(ips))

	for _, ip := range ips {

		// Skip IPv6 addresses for now.
		if ip.To4() == nil {
			continue
		}

		upstreams = append(upstreams, config.Upstream{
			URL:    fmt.Sprintf("%s://%s:%d", d.Scheme, ip.String(), d.Port),
			Weight: 1,
		})
	}

	return upstreams, nil
}
