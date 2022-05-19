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

type DnsEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Dns will only have its UID
	Dns object.Dns
}

func WatchDns(UID string) (chan DnsEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/dns/" + UID
	ch, cancel, err := createDnsWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

// WatchDnses
// if err != nil, chan and cancel() will be nil
// if you call cancel() or connection failed, channel will be closed
func WatchDnses() (chan DnsEvent, func(), error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/watch/dnses"
	ch, cancel, err := createDnsWatch(url)
	if err != nil && cancel != nil {
		cancelFuncs = append(cancelFuncs, cancel)
	}
	return ch, cancel, err
}

func createDnsWatch(url string) (chan DnsEvent, context.CancelFunc, error) {
	ch := make(chan DnsEvent)
	var closed int32 = 0
	closeChan := func() {
		swapped := atomic.CompareAndSwapInt32(&closed, 0, 1)
		if swapped {
			log.Println("closing DnsEvent channel")
			close(ch)
		}
	}
	stop, err := postWatch(url, closeChan, func(e ObjEvent) {
		var dnsEvent DnsEvent
		dnsEvent.EType = e.EType
		switch e.EType {
		case EVENT_PUT:
			err := json.Unmarshal([]byte(e.Object), &dnsEvent.Dns)
			if err != nil {
				log.Println("fail to parse Dns in DnsEvent")
				return
			}
		case EVENT_DELETE:
			dnsEvent.Dns.UID = e.Path[len(object.DnsEtcdPrefix):]
		}
		ch <- dnsEvent
	})
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		closeChan()
		stop()
	}, nil
}
