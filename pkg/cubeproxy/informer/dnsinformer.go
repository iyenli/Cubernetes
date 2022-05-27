package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

const WatchRetryIntervalSec = 10

type DNSInformer interface {
	ListAndWatchDNSWithRetry()
	WatchDNSEvent() <-chan types.DNSEvent
	ListDNS() []object.Dns
}

type ProxyDNSInformer struct {
	DNSChannel chan types.DNSEvent
	DNSCache   map[string]object.Dns

	mtx sync.RWMutex
}

func NewDNSInformer() DNSInformer {
	return &ProxyDNSInformer{
		DNSChannel: make(chan types.DNSEvent),
		DNSCache:   make(map[string]object.Dns),
	}
}

func (p *ProxyDNSInformer) ListAndWatchDNSWithRetry() {
	defer close(p.DNSChannel)
	for {
		p.tryListAndWatchDNS()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (p *ProxyDNSInformer) tryListAndWatchDNS() {
	if allDNS, err := crudobj.GetDnses(); err != nil {
		log.Printf("[Error]: fail to get all dnses from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		log.Printf("[INFO]: Ready to init dns using %v items\n", len(allDNS))
		for _, dns := range allDNS {
			p.DNSCache[dns.UID] = dns
		}
	}

	ch, cancel, err := watchobj.WatchDnses()
	if err != nil {
		log.Printf("[INFO]: fail to watch dnses from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case dnsEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			}
			switch dnsEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				err := p.informDNS(dnsEvent.Dns, dnsEvent.EType)
				if err != nil {
					log.Println("[Error]: Error when inform DNS:", dnsEvent.Dns.UID)
					return
				}
			default:
				log.Panic("[Fatal]: Unsupported types in watch dns")
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (p *ProxyDNSInformer) WatchDNSEvent() <-chan types.DNSEvent {
	return p.DNSChannel
}

func (p *ProxyDNSInformer) ListDNS() []object.Dns {
	p.mtx.RLock()
	dns := make([]object.Dns, len(p.DNSCache))
	idx := 0
	for _, item := range p.DNSCache {
		dns[idx] = item
		idx += 1
	}
	p.mtx.RUnlock()
	return dns
}

func (p *ProxyDNSInformer) informDNS(new object.Dns, eType watchobj.EventType) error {
	p.mtx.Lock()
	oldDns, exist := p.DNSCache[new.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(p.DNSCache, new.UID)
			p.DNSChannel <- types.DNSEvent{
				Type: types.Remove,
				DNS:  new,
			}
		} else {
			log.Printf("[INFO]: DNS %s not exist, delete do nothing\n", new.UID)
		}
	} else {
		// Any way update cache
		p.DNSCache[new.UID] = new
		if !exist {
			p.DNSChannel <- types.DNSEvent{
				Type: types.Create,
				DNS:  new,
			}
		} else {
			if object.ComputeDNSCriticalChange(&new, &oldDns) {
				log.Println("[INFO]: DNS critical change, UID is", new.UID)
				p.DNSChannel <- types.DNSEvent{
					Type: types.Update,
					DNS:  new,
				}
			} else {
				log.Println("[INFO]: DNS not change, UID is", new.UID)
			}
		}
	}
	p.mtx.Unlock()
	return nil
}
