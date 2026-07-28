package upstream

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func backend(raw string, active int64) *Backend {
	u, _ := url.Parse(raw)

	b := &Backend{
		URL:    u,
		Weight: 1,
	}

	b.Healthy.Store(true)
	b.CircuitState.Store(CircuitClosed)
	b.ActiveConnections.Store(active)

	return b
}

func TestLeastConnections_SelectsLeastBusyBackend(t *testing.T) {

	slow := backend("http://slow", 5)
	fast := backend("http://fast", 1)

	pool := &Pool{
		backends: []*Backend{
			slow,
			fast,
		},
		balancer: LeastConnectionsBalancer{},
	}

	selected := pool.Next()

	require.Same(t, fast, selected)
}

func TestLeastConnections_SelectsOtherBackend(t *testing.T) {

	slow := backend("http://slow", 1)
	fast := backend("http://fast", 8)

	pool := &Pool{
		backends: []*Backend{
			slow,
			fast,
		},
		balancer: LeastConnectionsBalancer{},
	}

	selected := pool.Next()

	require.Same(t, slow, selected)
}

func TestLeastConnections_SkipsUnhealthyBackend(t *testing.T) {

	b1 := backend("http://one", 0)
	b2 := backend("http://two", 0)

	b1.Healthy.Store(false)

	pool := &Pool{
		backends: []*Backend{
			b1,
			b2,
		},
		balancer: LeastConnectionsBalancer{},
	}

	selected := pool.Next()

	require.Same(t, b2, selected)
}

func TestLeastConnections_ReturnsNilWhenNoHealthyBackend(t *testing.T) {

	b1 := backend("http://one", 0)
	b2 := backend("http://two", 0)

	b1.Healthy.Store(false)
	b2.Healthy.Store(false)

	pool := &Pool{
		backends: []*Backend{
			b1,
			b2,
		},
		balancer: LeastConnectionsBalancer{},
	}

	require.Nil(t, pool.Next())
}
