package main

import (
	dockerapi "Cubernetes/pkg/cubelet/dockershim"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func main() {
	client, _ := dockerapi.NewDockerRuntime()
	defer client.CloseConnection()

	err := client.PullImage("busybox")
	if err != nil {
		panic(err)
	}

	id, err := client.CreateContainer(&types.ContainerCreateConfig{
		Name: "test-busybox",
		Config: &container.Config{
			Image: "busybox",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("id: ", id)
}