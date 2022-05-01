package loadbalancer

type LoadBalancer interface {
	InitBalancer() error
}
