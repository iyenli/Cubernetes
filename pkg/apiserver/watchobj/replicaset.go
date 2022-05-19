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

type ReplicaSetEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// ReplicaSet will only have its UID
	ReplicaSet object.ReplicaSet
}

func WatchReplicaSet(UID string) (chan ReplicaSetEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/replicaSet/" + UID
	ch, cancel, err := createReplicaSetWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchReplicaSets
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchReplicaSets() (chan ReplicaSetEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/replicaSets"
	ch, cancel, err := createReplicaSetWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createReplicaSetWatch(url string) (chan ReplicaSetEvent, context.CancelFunc, error) {
	ch := make(chan ReplicaSetEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing ReplicaSetEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var rsEvent ReplicaSetEvent
		rsEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &rsEvent.ReplicaSet)
			if err != nil {
				log.Println("fail to parse ReplicaSet in ReplicaSetEvent")
				return
			}
		case EVENT_DELETE:
			rsEvent.ReplicaSet.UID = e.Path[len(object.ReplicaSetEtcdPrefix):]
		}
		ch <- rsEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
