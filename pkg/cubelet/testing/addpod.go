package main

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func test() {
	time.Sleep(time.Second)
	for {
		targetUrl := "nginx:foo-bar-meaningless-uid"
		payload := strings.NewReader(`
{
  "kind": "Pod",
  "apiVersion": "wahtever/v1",
  "metadata": {
    "name": "hello",
    "uid": "nginx:foo-bar-meaningless-uid"
  },
  "spec": {
    "containers": [
      {
        "name": "foo-nginx",
        "image": "nginx",
        "ports": [
          {
            "containerPort": 8080,
            "hostPort": 80
          }
        ],
        "volumeMounts": [
          {
            "name": "conf",
            "mountPath": "/etc/nginx/nginx.conf"
          },
          {
            "name": "html",
            "mountPath": "/www"
          }
        ]
      }
    ],
    "volumes": [
      {
        "name": "conf",
        "hostPath": "/home/lee/CloudOS/test/nginx.conf/nginx.conf"
      },
      {
        "name": "html",
        "hostPath": "/home/lee/CloudOS/test/www"
      }
    ]
  }
}
`)
		req, _ := http.NewRequest("PUT", targetUrl, payload)
		req.Header.Add("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}

		log.Println(resp.Status)
		time.Sleep(15 * time.Second)
	}
}

// simple example of use
func main() {
	ch, cancel, _ := watchobj.WatchPods()
	go test()
	for podEvent := range ch {
		fmt.Println(podEvent)
	}
	// cancel watching
	cancel()
}
