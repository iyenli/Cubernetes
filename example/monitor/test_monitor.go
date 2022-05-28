package main

import (
	"Cubernetes/pkg/actionbrain/monitor"
	"fmt"
)

func main() {
	actionMonitor, err := monitor.NewActionMonitor("192.168.1.6")
	if err != nil {
		panic(err)
	}
	defer actionMonitor.Close()

	go actionMonitor.Run()
	ch := actionMonitor.WatchActionEvoke()

	for action := range ch {
		fmt.Printf("receive Action: %s\n", action)

		if action == "quit" {
			fmt.Printf("quit monitor test...\n")
		}
	}
}
