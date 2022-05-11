package main

import (
	"Cubernetes/pkg/cubelet"
	"Cubernetes/pkg/cubelet/network"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"log"
	"os"
)

func main() {
	// Init network according to params
	// 1 param (master): cubelet [LocalIP]
	// 2 params (slave): cubelet [LocalIP] [MasterIP]
	if len(os.Args) < 2 {
		log.Fatal("[FATAL] Lack arguments")
	}

	if len(os.Args) == 3 {
		nodenetwork.SetMasterIP(os.Args[2])
	}

	network.InitNodeNetwork(os.Args)
	network.InitNodeHeartbeat()
	cubeletInstance := cubelet.NewCubelet()
	cubeletInstance.Run()
}
