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

func MatchServiceAndPod(service *object.Service, pod *object.Pod) bool {
	for sLabelKey, sLabelVal := range service.Labels {
		hasThisLabel := false
		for pLabelKey, pLabelVal := range pod.Labels {
			if pLabelKey == sLabelKey {
				if pLabelVal == sLabelVal {
					hasThisLabel = true
					break
				} else {
					return false
				}
			}
		}

		if !hasThisLabel {
			return false
		}
	}

	return true
}
