package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
)

type DNSInformer interface {
	InitInformer(dns []object.Dns) error
	WatchDNSEvent() <-chan types.DNSEvent
	InformDNS(new object.Dns, eType watchobj.EventType) error
	ListDNS() []object.Dns
	CloseChan()
}

type ProxyDNSInformer struct {
	DNSChannel chan types.DNSEvent
	DNSCache   map[string]object.Dns
}

func NewDNSInformer() DNSInformer {
	return &ProxyDNSInformer{
		DNSChannel: make(chan types.DNSEvent),
		DNSCache:   make(map[string]object.Dns),
	}
}

func (p *ProxyDNSInformer) InitInformer(dns []object.Dns) error {
	for _, item := range dns {
		p.DNSCache[item.UID] = item
	}

	return nil
}

func (p *ProxyDNSInformer) WatchDNSEvent() <-chan types.DNSEvent {
	return p.DNSChannel
}

func (p *ProxyDNSInformer) CloseChan() {
	close(p.DNSChannel)
}

func (p *ProxyDNSInformer) ListDNS() []object.Dns {
	dns := make([]object.Dns, len(p.DNSCache))
	idx := 0
	for _, item := range p.DNSCache {
		dns[idx] = item
		idx += 1
	}

	return dns
}

func (p *ProxyDNSInformer) InformDNS(new object.Dns, eType watchobj.EventType) error {
	oldDns, exist := p.DNSCache[new.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(p.DNSCache, new.UID)
			p.DNSChannel <- types.DNSEvent{
				Type: types.DNSRemove,
				DNS:  new,
			}
		} else {
			log.Printf("[INFO]: pod %s not exist, delete do nothing\n", new.UID)
		}
	} else {
		// Any way update cache
		p.DNSCache[new.UID] = new
		if !exist {
			p.DNSChannel <- types.DNSEvent{
				Type: types.DNSCreate,
				DNS:  new,
			}
		} else {
			if object.ComputeDNSCriticalChange(&new, &oldDns) {
				log.Println("[INFO]: DNS critical change, UID is", new.UID)
				p.DNSChannel <- types.DNSEvent{
					Type: types.DNSUpdate,
					DNS:  new,
				}
			} else {
				log.Println("[INFO]: DNS not change, UID is", new.UID)
			}
		}
	}

	return nil
}
