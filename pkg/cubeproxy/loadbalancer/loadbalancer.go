package loadbalancer

// LoadBalancer TODO: More complicated load balancer
type LoadBalancer interface {
	InitBalancer() error
}
