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

type GpuJobEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// GpuJob will only have its UID
	GpuJob object.GpuJob
}

func WatchGpuJob(UID string) (chan GpuJobEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/gpuJob/" + UID
	ch, cancel, err := createGpuJobWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchGpuJobs
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchGpuJobs() (chan GpuJobEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/gpuJobs"
	ch, cancel, err := createGpuJobWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createGpuJobWatch(url string) (chan GpuJobEvent, context.CancelFunc, error) {
	ch := make(chan GpuJobEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing GpuJobEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var gpuJobEvent GpuJobEvent
		gpuJobEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &gpuJobEvent.GpuJob)
			if err != nil {
				log.Println("fail to parse GpuJob in GpuJobEvent")
				return
			}
		case EVENT_DELETE:
			gpuJobEvent.GpuJob.UID = e.Path[len(object.GpuJobEtcdPrefix):]
		}
		ch <- gpuJobEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
