package main

import (
	"Cubernetes/pkg/watchobj"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func test() {
	time.Sleep(time.Second)
	for {
		targetUrl := "http://127.0.0.1:8080/apis/pod/hello:e0a77a11-f736-4f5f-934e-f1f0a3c39172"
		payload := strings.NewReader("{\"kind\":\"pod\",\"apiVersion\":\"3\",\"metadata\":{\"name\":\"hello\",\"uid\":\"hello:e0a77a11-f736-4f5f-934e-f1f0a3c39172\"},\"spec\":{\"containers\":null},\"status\":{}}")
		req, _ := http.NewRequest("PUT", targetUrl, payload)
		req.Header.Add("Content-Type", "application/json")
		_, _ = http.DefaultClient.Do(req)
		time.Sleep(time.Millisecond)
	}
}

func waitAndCancel(cancel func()) {
	time.Sleep(5 * time.Second)
	cancel()
}

// simple example of use
func main() {
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		fmt.Println("error")
		return
	}
	//go test()
	// example for cancel watching
	go waitAndCancel(cancel)

	fmt.Println("start watching")
	for podEvent := range ch {
		fmt.Println(podEvent)
	}
	fmt.Println("watch cancelled")
}
