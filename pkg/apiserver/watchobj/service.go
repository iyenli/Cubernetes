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

type ServiceEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Service will only have its UID
	Service object.Service
}

func WatchService(UID string) (chan ServiceEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/service/" + UID
	ch, cancel, err := createServiceWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchServices
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchServices() (chan ServiceEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/services"
	ch, cancel, err := createServiceWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createServiceWatch(url string) (chan ServiceEvent, context.CancelFunc, error) {
	ch := make(chan ServiceEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing ServiceEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var serviceEvent ServiceEvent
		serviceEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &serviceEvent.Service)
			if err != nil {
				log.Println("fail to parse Service in ServiceEvent")
				return
			}
		case EVENT_DELETE:
			serviceEvent.Service.UID = e.Path[len(object.ServiceEtcdPrefix):]
		}
		ch <- serviceEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
