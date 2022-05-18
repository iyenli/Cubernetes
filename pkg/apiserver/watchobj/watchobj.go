package watchobj

import (
	cubeconfig "Cubernetes/config"
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type EventType string

const (
	EVENT_PUT    EventType = "PUT"
	EVENT_DELETE EventType = "DELETE"
)

const MSG_DELIM byte = 26
const WATCH_CONFIRM string = "watch started"

type ObjEvent struct {
	EType  EventType `json:"eType"`
	Path   string    `json:"path"`
	Object string    `json:"object"`
}

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
