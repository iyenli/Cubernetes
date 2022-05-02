package proxyruntime

import (
	"Cubernetes/pkg/object"
	"net"
)

// GetPodByService TODO: Change to real pod interface
func GetPodByService(service *object.Service) ([]object.Pod, error) {
	//return crudobj.SelectPods(service.Spec.Selector)
	pods := []object.Pod{
		{
			TypeMeta: object.TypeMeta{
				Kind:       "Pod",
				APIVersion: "c8s/v1",
			},
			ObjectMeta: object.ObjectMeta{
				Name:        "superPod",
				Namespace:   "ns",
				UID:         "fake",
				Labels:      nil,
				Annotations: nil,
			},
			Spec: object.PodSpec{
				Containers: []object.Container{
					{
						Name:         "nginx",
						Image:        "nginx",
						Command:      nil,
						Args:         nil,
						Resources:    nil,
						VolumeMounts: nil,
						Ports: []object.ContainerPort{
							{
								Name:          "port",
								HostPort:      7895,
								ContainerPort: 1234,
								Protocol:      "TCP",
								HostIP:        "127.0.0.1",
							},
						},
					},
				},
				Volumes: nil,
			},
			Status: &object.PodStatus{
				IP: net.ParseIP("10.0.0.12"),
			},
		},
	}

	return pods, nil
}
