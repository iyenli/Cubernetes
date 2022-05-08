package network

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"log"
	"net"
)

func InitNodeNetwork(args []string) {
	var err error
	if len(args) == 2 {
		// master
		err = weaveplugins.InitWeave()
	} else if len(args) == 3 {
		// slave
		err = weaveplugins.AddNode(weaveplugins.Host{IP: net.ParseIP(args[1])}, weaveplugins.Host{IP: net.ParseIP(args[2])})
		registerToAPIServer(args)
		cubeconfig.APIServerIp = args[2]
	} else {
		panic("Error: too much or little args when start cubelet;")
	}

	if err != nil {
		log.Panicf("Init weave network failed, err: %v", err.Error())
		return
	}
}

func registerToAPIServer(args []string) {
	// TODO: Register to API Server
}
