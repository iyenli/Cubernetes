package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
)

// GetPodByService
// For test purpose:)
func GetPodByService(service *object.Service) ([]object.Pod, error) {
	//return crudobj.SelectPods(service.Spec.Selector)
	pods, err := crudobj.SelectPods(service.Spec.Selector)
	if err != nil {
		log.Println("[ERROR]: failed when select pods by service")
		return nil, err
	}

	return pods, nil
}
