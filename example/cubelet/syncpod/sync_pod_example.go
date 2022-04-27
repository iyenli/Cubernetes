package main

import (
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/cubelet/cuberuntime"
	"Cubernetes/pkg/cubelet/container"
)

func main() {
	pod := &object.Pod{
		TypeMeta: object.TypeMeta{
			Kind: "Pod",
			APIVersion: "wahtever/v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "test-pod",
		},
		Spec: object.PodSpec{
			Containers: []object.Container{
				{
					Name: "foo-nginx",
					Image: "nginx",
					Ports: []object.ContainerPort{
						{
							HostPort: 8080,
							ContainerPort: 80,
						},
					},
					VolumeMounts: []object.VolumeMount{
						{
							Name: "conf",
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

	runtime, err := cuberuntime.NewCubeRuntimeManager()
	if err != nil {
		panic(err)
	}
	defer runtime.Close()

	err = runtime.SyncPod(pod, &container.PodStatus{})
	if err != nil {
		panic(err)
	}

}