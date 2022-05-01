package main

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"fmt"
	"time"
)

func test() {
	time.Sleep(time.Second)
	for {
		var pod object.Pod
		pod.APIVersion = "1"
		pod.Kind = "pod"
		pod.Name = "hello233"

		pod, err := crudobj.CreatePod(pod)
		if err != nil {
			return
		}

		time.Sleep(200 * time.Millisecond)

		err = crudobj.DeletePod(pod.UID)
		if err != nil {
			return
		}

		time.Sleep(200 * time.Millisecond)
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
	go test()
	// example for cancel watching
	go waitAndCancel(cancel)

	fmt.Println("start watching")
	for podEvent := range ch {
		fmt.Println(podEvent)
	}
	fmt.Println("watch cancelled")
}
