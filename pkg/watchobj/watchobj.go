package watchobj

import (
	cubeconfig "Cubernetes/config"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func watching(ctx context.Context, url string, handler func(ObjEvent)) {
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	data := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readN, err := resp.Body.Read(data)
			if readN > 0 {
				var objEvent ObjEvent
				err := json.Unmarshal(data[:readN], &objEvent)
				if err != nil {
					log.Println("fail to parse objEvent")
					continue
				}
				handler(objEvent)
			}
			if err == io.EOF {
				log.Println("connection closed by server")
				return
			}
			if err != nil {
				panic(err)
			}
		}
	}
}

func WatchObj(path string) (chan ObjEvent, context.CancelFunc) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + path
	ch := make(chan ObjEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	go watching(ctx, url, func(e ObjEvent) {
		ch <- e
	})
	return ch, cancel
}

func createPodWatch(url string) (chan PodEvent, context.CancelFunc) {
	ch := make(chan PodEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	go watching(ctx, url, func(e ObjEvent) {
		var podEvent PodEvent
		podEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &podEvent.Pod)
			if err != nil {
				log.Println("fail to parse Pod in PodEvent")
				return
			}
		case EVENT_DELETE:
			podEvent.Pod.UID = e.Path[len("/apis/pod/"):]
		}
		ch <- podEvent
	})
	return ch, cancel
}

func WatchPod(UID string) (chan PodEvent, context.CancelFunc) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pod/" + UID
	return createPodWatch(url)
}

func WatchPods() (chan PodEvent, context.CancelFunc) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pods"
	return createPodWatch(url)
}
