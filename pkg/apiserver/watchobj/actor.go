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

type ActorEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Actor will only have its UID
	Actor object.Actor
}

func WatchActor(UID string) (chan ActorEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/actor/" + UID
	ch, cancel, err := createActorWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchActors
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchActors() (chan ActorEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/actors"
	ch, cancel, err := createActorWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createActorWatch(url string) (chan ActorEvent, context.CancelFunc, error) {
	ch := make(chan ActorEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing ActorEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var actorEvent ActorEvent
		actorEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &actorEvent.Actor)
			if err != nil {
				log.Println("fail to parse Actor in ActorEvent")
				return
			}
		case EVENT_DELETE:
			actorEvent.Actor.UID = e.Path[len(object.ActorEtcdPrefix):]
		}
		ch <- actorEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
