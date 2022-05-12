package main

import (
	"Cubernetes/cmd/apiserver/heartbeat"
	"Cubernetes/cmd/apiserver/httpserver"
	"Cubernetes/cmd/apiserver/httpserver/restful"
	"Cubernetes/pkg/cubenetwork/servicenetwork"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		go heartbeat.ListenHeartbeat()
	}()
	go func() {
		defer wg.Done()
		httpserver.Run()
	}()

	time.Sleep(time.Second)
	restful.ClusterIPAllocator = servicenetwork.NewClusterIPAllocator()

	log.Println("Cluster IP Allocator init, api server running...")
	wg.Wait()
}
