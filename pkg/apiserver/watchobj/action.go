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

type ActionEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Action will only have its UID
	Action object.Action
}

func WatchAction(UID string) (chan ActionEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/action/" + UID
	ch, cancel, err := createActionWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchActions
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchActions() (chan ActionEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/actions"
	ch, cancel, err := createActionWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createActionWatch(url string) (chan ActionEvent, context.CancelFunc, error) {
	ch := make(chan ActionEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing ActionEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var actionEvent ActionEvent
		actionEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &actionEvent.Action)
			if err != nil {
				log.Println("fail to parse Action in ActionEvent")
				return
			}
		case EVENT_DELETE:
			actionEvent.Action.UID = e.Path[len(object.ActionEtcdPrefix):]
		}
		ch <- actionEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
