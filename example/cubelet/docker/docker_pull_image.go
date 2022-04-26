package main

import (dockerapi "Cubernetes/pkg/cubelet/dockershim")

func main() {
	client, _ := dockerapi.NewDockerRuntime()
	err := client.PullImage("busybox")
	if err != nil {
		panic(err)
	}
	client.CloseConnection()
}