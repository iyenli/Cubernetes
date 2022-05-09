package main

import (
	"Cubernetes/pkg/cubenetwork/register"
	"Cubernetes/pkg/cubeproxy"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"os"
)

func main() {
	runtime, err := proxyruntime.InitIPTables()
	if err != nil {
		log.Printf("Create cube proxy runtime error: %v", err.Error())
	}

	register.RegisterToMaster(os.Args)
	cubeProxyInstance := cubeproxy.Cubeproxy{Runtime: runtime}
	cubeProxyInstance.Run()
}
