package main

import (
	"Cubernetes/cmd/apiserver/heartbeat"
	"Cubernetes/cmd/apiserver/httpserver"
)

func main() {
	go heartbeat.ListenHeartbeat()
	httpserver.Run()
}
