package upstream

type Balancer interface {
	Next() *Backend
	HasHealthyBackend() bool
}
