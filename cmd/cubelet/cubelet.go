package main

import (
	"Cubernetes/pkg/cubelet"
	"Cubernetes/pkg/cubelet/network"
	"os"
)

func main() {
	// Init network according to params, Register to api server
	network.InitNodeNetwork(os.Args)
	cubeletInstance := cubelet.NewCubelet()
	cubeletInstance.Run()
}
