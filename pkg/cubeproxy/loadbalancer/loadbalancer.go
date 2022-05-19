package loadbalancer

import "Cubernetes/pkg/object"

// LoadBalancer TODO: More complicated load balancer
type LoadBalancer interface {
	InitBalancer() error

	NextPodIP() string

	AddPod(pod *object.Pod) error
}

type RRLB struct {
	queue []int
}

func (*RRLB) InitBalancer() error {
	return nil
}

func (*RRLB) NextPodIP() string {
	return ""
}

func (*RRLB) AddPod(pod *object.Pod) error {
	return nil
}
