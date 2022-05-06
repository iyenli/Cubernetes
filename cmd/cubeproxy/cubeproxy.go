package main

import (
	"Cubernetes/pkg/cubeproxy"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
)

func main() {
	runtime, err := proxyruntime.InitIPTables()
	if err != nil {
		log.Printf("Create cube proxy runtime error: %v", err.Error())
	}

	cubeProxyInstance := cubeproxy.Cubeproxy{Runtime: runtime}
	cubeProxyInstance.Run()
}
