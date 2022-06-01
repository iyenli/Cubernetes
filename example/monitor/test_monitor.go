package main

import (
	"Cubernetes/pkg/actionbrain/monitor"
	"fmt"
	"log"
	"time"
)

func listen(ch <-chan string) {
	for action := range ch {
		fmt.Printf("receive action %s\n", action)
	}
}

func main() {
	actionMonitor, err := monitor.NewActionMonitor("192.168.1.6")
	if err != nil {
		panic(err)
	}
	defer actionMonitor.Close()

	go actionMonitor.Run()
	
	log.Printf("you can start to send")
	go listen(actionMonitor.WatchActionEvoke())
	time.Sleep(time.Second * 40)

	count, err := actionMonitor.QueryRecentEvoke("fuck-me", time.Hour)
	if err != nil {
		panic(err)
	}

	log.Printf("total %d records counted", count)
}
