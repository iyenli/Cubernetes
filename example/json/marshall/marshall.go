package main

import (
	"Cubernetes/pkg/object"
	"fmt"

	"encoding/json"
)

func main() {

	pod := &object.Pod{
		TypeMeta: object.TypeMeta{
			Kind:       "Pod",
			APIVersion: "wahtever/v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "test-pod",
		},
		Spec: object.PodSpec{
			Containers: []object.Container{
				{
					Name:  "foo-nginx",
					Image: "nginx",
					Ports: []object.ContainerPort{
						{
							HostPort:      8080,
							ContainerPort: 80,
						},
					},
					VolumeMounts: []object.VolumeMount{
						{
							Name:      "conf",
							MountPath: "/etc/nginx/nginx.conf",
						},
					},
				},
			},
			Volumes: []object.Volume{
				{
					Name: "conf",
					// switch to your host-path when test
					HostPath: "/home/jolynefr/WorkSpace/CloudOS/test/nginx.conf/nginx.conf",
				},
			},
		},
	}

	b, err := json.Marshal(pod)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}