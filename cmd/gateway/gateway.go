package main

import (
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/gateway"
	"log"
	"os"
)

func main() {
	// Init network according to params
	// 1 param: cubelet [MasterIP]
	if len(os.Args) != 2 {
		log.Fatal("[FATAL] Lack arguments")
	}

	nodenetwork.SetMasterIP(os.Args[1])
	runtime := gateway.NewRuntimeGateway()
	if runtime == nil {
		log.Panicln("[Error]: init gateway failed")
	}

	runtime.Run()
	log.Panicln("[Fatal]: unreachable here")
}
