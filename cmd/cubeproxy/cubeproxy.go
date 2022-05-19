package main

import (
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubeproxy"
	"log"
	"os"
)

func main() {
	// 1 param (master): cubeproxy [LocalIP]
	// 2 params (slave): cubeproxy [LocalIP] [MasterIP]
	if len(os.Args) < 2 {
		log.Fatal("[FATAL] Lack arguments")
	}
	if len(os.Args) == 3 {
		nodenetwork.SetMasterIP(os.Args[2])
	}

	cubeProxyInstance := cubeproxy.NewCubeProxy()
	cubeProxyInstance.Run()

	log.Fatalln("[Fatal]: Unreachable here")
}
