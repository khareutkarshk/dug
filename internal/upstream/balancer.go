package upstream

type Balancer interface {
	Next(*Pool) *Backend
}
