package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
)

func GetPodByService(service *object.Service) ([]object.Pod, error) {
	return crudobj.SelectPods(service.Spec.Selector)
}
