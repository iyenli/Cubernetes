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

type NodeEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Node will only have its UID
	Node object.Node
}

func WatchNode(UID string) (chan NodeEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/node/" + UID
	ch, cancel, err := createNodeWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchNodes
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchNodes() (chan NodeEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/nodes"
	ch, cancel, err := createNodeWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createNodeWatch(url string) (chan NodeEvent, context.CancelFunc, error) {
	ch := make(chan NodeEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing NodeEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var nodeEvent NodeEvent
		nodeEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &nodeEvent.Node)
			if err != nil {
				log.Println("fail to parse Node in NodeEvent")
				return
			}
		case EVENT_DELETE:
			nodeEvent.Node.UID = e.Path[len(object.NodeEtcdPrefix):]
		}
		ch <- nodeEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
