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

type IngressEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Ingress will only have its UID
	Ingress object.Ingress
}

func WatchIngress(UID string) (chan IngressEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/ingress/" + UID
	ch, cancel, err := createIngressWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchIngresses
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchIngresses() (chan IngressEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/ingresses"
	ch, cancel, err := createIngressWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createIngressWatch(url string) (chan IngressEvent, context.CancelFunc, error) {
	ch := make(chan IngressEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing IngressEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var ingressEvent IngressEvent
		ingressEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &ingressEvent.Ingress)
			if err != nil {
				log.Println("fail to parse Ingress in IngressEvent")
				return
			}
		case EVENT_DELETE:
			ingressEvent.Ingress.UID = e.Path[len(object.IngressEtcdPrefix):]
		}
		ch <- ingressEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
