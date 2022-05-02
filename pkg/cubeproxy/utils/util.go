package utils

import "Cubernetes/pkg/object"

// DefaultService
// 1. Check if service is legal
// 2. Make it legal using default values in K8
func DefaultService(service *object.Service) error {
	for _, port := range service.Spec.Ports {
		if port.Port == 0 {
			port.Port = port.TargetPort
		}
	}
	return nil
}
