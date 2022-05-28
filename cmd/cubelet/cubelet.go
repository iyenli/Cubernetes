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
	// 2 param (master): cubelet [UID Of Node] [LocalIP]
	// 3 params (slave): cubelet [UID Of Node] [LocalIP] [MasterIP]
	if len(os.Args) < 3 {
		log.Fatal("[FATAL] Lack arguments")
	}

	if len(os.Args) == 3 {
		nodenetwork.SetMasterIP(os.Args[2])
	} else if len(os.Args) == 4 {
		nodenetwork.SetMasterIP(os.Args[3])
	}

	ip := network.InitNodeNetwork(os.Args)
	network.InitNodeHeartbeat()

	cubeletInstance := cubelet.NewCubelet()
	cubeletInstance.InitCubelet(os.Args[1], ip)
	cubeletInstance.Run()
}
