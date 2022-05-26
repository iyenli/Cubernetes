package main

import (
	"Cubernetes/pkg/gateway"
)

func main() {
	// Init network according to params
	// 1 param: cubelet [MasterIP]
	// TODO: fill in master IP
	runtime := gateway.NewRuntimeGateway()
	runtime.Run()
}
