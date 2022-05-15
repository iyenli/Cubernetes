package watchobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"sync/atomic"
)

type PodEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Pod will only have its UID
	Pod object.Pod
}

func WatchPod(UID string) (chan PodEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pod/" + UID
	ch, cancel, err := createPodWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchPods
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchPods() (chan PodEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pods"
	ch, cancel, err := createPodWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createPodWatch(url string) (chan PodEvent, context.CancelFunc, error) {
	ch := make(chan PodEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing PodEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
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
			podEvent.Pod.UID = e.Path[len(object.PodEtcdPrefix):]
		}
		ch <- podEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
