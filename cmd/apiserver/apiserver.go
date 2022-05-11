package main

import (
	"Cubernetes/cmd/apiserver/heartbeat"
	"Cubernetes/cmd/apiserver/httpserver"
	"Cubernetes/pkg/utils/etcdrw"
)

func main() {
	etcdrw.Init()
	defer etcdrw.Free()

	go heartbeat.ListenHeartbeat()
	httpserver.Run()
}
