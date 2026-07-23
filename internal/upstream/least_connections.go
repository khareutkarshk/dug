package upstream

type LeastConnectionsBalancer struct{}

func (LeastConnectionsBalancer) Next(p *Pool) *Backend {
	panic("Least connection is not implemented yet")
}
