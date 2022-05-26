package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/gateway/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

const WatchRetryIntervalSec = 10

type IngressInformer interface {
	ListAndWatchIngressWithRetry()
	WatchIngressEvent() <-chan types.IngressEvent
	ListIngress() []object.Ingress
	GetIngressLocally(UID string) *object.Ingress
}

func NewIngressInformer() IngressInformer {
	return &ProxyIngressInformer{
		IngressChannel: make(chan types.IngressEvent),
		IngressCache:   make(map[string]object.Ingress),
	}
}

type ProxyIngressInformer struct {
	IngressChannel chan types.IngressEvent
	IngressCache   map[string]object.Ingress
	PathSet        map[string]struct{}

	mtx sync.RWMutex
}

func (p *ProxyIngressInformer) GetIngressLocally(UID string) *object.Ingress {
	if ingress, ok := p.IngressCache[UID]; ok {
		return &ingress
	}
	return nil
}

func (p *ProxyIngressInformer) ListAndWatchIngressWithRetry() {
	defer close(p.IngressChannel)
	for {
		p.tryListAndWatchIngress()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (p *ProxyIngressInformer) tryListAndWatchIngress() {
	if allIngress, err := crudobj.GetIngresses(); err != nil {
		log.Printf("[Error]: fail to get all Ingresses from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		log.Printf("[INFO]: Ready to init Ingress using %v items\n", len(allIngress))
		for _, Ingress := range allIngress {
			p.IngressCache[Ingress.UID] = Ingress
		}
	}

	ch, cancel, err := watchobj.WatchIngresses()
	if err != nil {
		log.Printf("[INFO]: fail to watch Ingresses in Gateway: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case IngressEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			}
			switch IngressEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				err := p.informIngress(IngressEvent.Ingress, IngressEvent.EType)
				if err != nil {
					log.Println("[Error]: Error when inform Ingress:", IngressEvent.Ingress.UID)
					return
				}
			default:
				log.Panic("[Fatal]: Unsupported types in watch Ingress")
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (p *ProxyIngressInformer) WatchIngressEvent() <-chan types.IngressEvent {
	return p.IngressChannel
}

func (p *ProxyIngressInformer) ListIngress() []object.Ingress {
	p.mtx.RLock()

	Ingress := make([]object.Ingress, len(p.IngressCache))
	idx := 0
	for _, item := range p.IngressCache {
		Ingress[idx] = item
		idx += 1
	}

	p.mtx.RUnlock()
	return Ingress
}

func (p *ProxyIngressInformer) informIngress(new object.Ingress, eType watchobj.EventType) error {
	p.mtx.Lock()
	oldIngress, exist := p.IngressCache[new.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			p.IngressChannel <- types.IngressEvent{
				Type:    types.IngressRemove,
				Ingress: oldIngress,
			}
			delete(p.IngressCache, new.UID)
		} else {
			log.Printf("[INFO]: Ingress %s not exist, delete do nothing\n", new.UID)
		}
	} else {
		// Duplicate Ingress trigger path detect
		if !exist {
			if _, ok := p.PathSet[new.Spec.TriggerPath]; ok {
				log.Printf("[Warn]: UID is different but trigger path duplicate")
				return nil
			}
		}

		// Any way update cache
		p.IngressCache[new.UID] = new
		if !exist {
			p.IngressChannel <- types.IngressEvent{
				Type:    types.IngressCreate,
				Ingress: new,
			}
		} else {
			if object.ComputeIngressCriticalChange(&new, &oldIngress) {
				log.Println("[INFO]: Ingress critical change, UID is", new.UID)
				p.IngressChannel <- types.IngressEvent{
					Type:    types.IngressUpdate,
					Ingress: new,
				}
			} else {
				log.Println("[INFO]: Ingress not change, UID is", new.UID)
			}
		}
	}

	p.mtx.Unlock()
	return nil
}
