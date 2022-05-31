package main

import (
	"Cubernetes/pkg/actionbrain/monitor"
	"log"
	"time"
)

func main() {
	actionMonitor, err := monitor.NewActionMonitor("192.168.1.6")
	if err != nil {
		panic(err)
	}
	defer actionMonitor.Close()

	go actionMonitor.Run()
	
	log.Printf("you can start to send")
	time.Sleep(time.Second * 40)

	count, err := actionMonitor.QueryRecentEvoke("fuck-you", time.Hour)
	if err != nil {
		panic(err)
	}

	log.Printf("total %d records counted", count)
}
