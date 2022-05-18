package testing

import (
	"Cubernetes/pkg/object"
	"net"
	"testing"
)

var pods = []object.Pod{
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
							// Host machine IP
							HostIP: "172.0.0.3",
						},
					},
				},
			},
			Volumes: nil,
		},
		Status: &object.PodStatus{
			// Pod IP
			IP: net.ParseIP("10.0.0.12"),
		},
	},
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
							// Host machine IP
							HostIP: "172.0.0.4",
						},
					},
				},
			},
			Volumes: nil,
		},
		Status: &object.PodStatus{
			// Pod IP
			IP: net.ParseIP("10.0.0.15"),
		},
	},
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
							// Host machine IP
							HostIP: "172.0.0.3",
						},
					},
				},
			},
			Volumes: nil,
		},
		Status: &object.PodStatus{
			// Pod IP
			IP: net.ParseIP("10.0.0.14"),
		},
	},
}

func TestCubeproxy(t *testing.T) {

}
