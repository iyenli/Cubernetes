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

type AutoScalerEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// AutoScaler will only have its UID
	AutoScaler object.AutoScaler
}

func WatchAutoScaler(UID string) (chan AutoScalerEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/autoScaler/" + UID
	ch, cancel, err := createAutoScalerWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchAutoScalers
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchAutoScalers() (chan AutoScalerEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/autoScalers"
	ch, cancel, err := createAutoScalerWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createAutoScalerWatch(url string) (chan AutoScalerEvent, context.CancelFunc, error) {
	ch := make(chan AutoScalerEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing AutoScalerEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var asEvent AutoScalerEvent
		asEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &asEvent.AutoScaler)
			if err != nil {
				log.Println("fail to parse AutoScaler in AutoScalerEvent")
				return
			}
		case EVENT_DELETE:
			asEvent.AutoScaler.UID = e.Path[len(object.AutoScalerEtcdPrefix):]
		}
		ch <- asEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
