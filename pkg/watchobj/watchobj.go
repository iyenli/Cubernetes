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
	return ch, func() {
		closeChan()
		stop()
	}, nil
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

func WatchPod(UID string) (chan PodEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pod/" + UID
	return createPodWatch(url)
}

// WatchPods
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchPods() (chan PodEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/pods"
	return createPodWatch(url)
}
