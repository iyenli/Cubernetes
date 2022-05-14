package watchobj

import (
	cubeconfig "Cubernetes/config"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var cancelFuncs []func()

func StopAll() {
	for _, cancel := range cancelFuncs {
		if cancel != nil {
			cancel()
		}
	}
	cancelFuncs = cancelFuncs[0:0]
}

func watching(reader *bufio.Reader, closeChan func(), handler func(ObjEvent)) {
	defer closeChan()
	for {
		buf, err := reader.ReadBytes(MSG_DELIM)
		if err != nil {
			log.Println("connection closed")
			return
		}
		var objEvent ObjEvent
		err = json.Unmarshal(buf[:len(buf)-1], &objEvent)
		if err != nil {
			log.Println("fail to parse objEvent")
			continue
		}
		handler(objEvent)
	}
}

func postWatch(url string, closeChan func(), handler func(ObjEvent)) (func(), error) {
	resp, err := http.Post(url, "application/json", strings.NewReader("{}"))
	if err != nil {
		log.Println("fail to send http post request")
		return nil, err
	}
	reader := bufio.NewReader(resp.Body)
	buf, err := reader.ReadBytes(MSG_DELIM)
	if err != nil {
		log.Println("connection closed")
		_ = resp.Body.Close()
		return nil, err
	}
	if string(buf[:len(buf)-1]) != WATCH_CONFIRM {
		log.Println("server did not confirm watching")
		_ = resp.Body.Close()
		return nil, errors.New("server did not confirm watching")
	}
	go watching(reader, closeChan, handler)
	return func() { _ = resp.Body.Close() }, nil
}

func WatchObj(path string) (chan ObjEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + path
	ch := make(chan ObjEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing ObjEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		ch <- e
	})
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		closeChan()
		stop()
	}
	cancelFuncs = append(cancelFuncs, cancel)
	return ch, cancel, nil
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
			podEvent.Pod.UID = e.Path[len("/apis/pod/"):]
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
			serviceEvent.Service.UID = e.Path[len("/apis/service/"):]
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
			rsEvent.ReplicaSet.UID = e.Path[len("/apis/replicaSet/"):]
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
				log.Println("fail to parse Node in ReplicaSetEvent")
				return
			}
		case EVENT_DELETE:
			nodeEvent.Node.UID = e.Path[len("/apis/node/"):]
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
