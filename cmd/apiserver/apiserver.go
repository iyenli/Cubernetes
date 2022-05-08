package main

import "Cubernetes/cmd/apiserver/httpserver"

func main() {
	go listenHeartbeat()
	httpserver.Run()
}
