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

// simple example of use
func main() {
	ch, cancel := watchobj.WatchObj("/apis/watch/pod/hello:e0a77a11-f736-4f5f-934e-f1f0a3c39172")
	go test()
	for str := range ch {
		fmt.Println(str)
	}
	// cancel watching
	cancel()
}
