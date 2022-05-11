package main

import (
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubeproxy"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"os"
)

func main() {
	// 1 param (master): cubeproxy [LocalIP]
	// 2 params (slave): cubeproxy [LocalIP] [MasterIP]
	if len(os.Args) < 2 {
		log.Fatal("[FATAL] Lack arguments")
	}

	runtime, err := proxyruntime.InitIPTables()
	if err != nil {
		log.Printf("Create cube proxy runtime error: %v", err.Error())
	}

	if len(os.Args) == 3 {
		nodenetwork.SetMasterIP(os.Args[2])
	}
	cubeProxyInstance := cubeproxy.Cubeproxy{Runtime: runtime}
	cubeProxyInstance.Run()
}
