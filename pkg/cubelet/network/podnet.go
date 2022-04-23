package network

import (
	"Cubernetes/pkg/object"
	"net"
)

// ConstructPodPortMapping creates a PodPortMapping from the ports specified in the pod's
// containers.
func ConstructPodPortMapping(pod *object.Pod, podIP net.IP) *PodPortMapping {
	portMappings := make([]*PortMapping, 0)
	for _, c := range pod.Spec.Containers {
		for _, port := range c.Ports {
			portMappings = append(portMappings, &PortMapping{
				Name:          port.Name,
				HostPort:      port.HostPort,
				ContainerPort: port.ContainerPort,
				Protocol:      port.Protocol,
				HostIP:        port.HostIP,
			})
		}
	}

	return &PodPortMapping{
		Namespace:    pod.Namespace,
		Name:         pod.Name,
		PortMappings: portMappings,
		IP:           podIP,
	}
}
