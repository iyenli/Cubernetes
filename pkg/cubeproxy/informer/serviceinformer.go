package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
)

type ServiceInformer interface {
	InitInformer(services []object.Service) error
	WatchServiceEvent() <-chan types.ServiceEvent
	InformService(newService object.Service, eType watchobj.EventType) error
	ListServices() []object.Service
	CloseChan()
}

type ProxyServiceInformer struct {
	ServiceChannel chan types.ServiceEvent
	ServiceCache   map[string]object.Service

	mtx sync.RWMutex
}

func NewServiceInformer() ServiceInformer {
	return &ProxyServiceInformer{
		ServiceChannel: make(chan types.ServiceEvent),
		ServiceCache:   make(map[string]object.Service),
	}
}

func (i *ProxyServiceInformer) InitInformer(services []object.Service) error {
	for _, service := range services {
		i.ServiceCache[service.UID] = service
	}

	return nil
}

func (i *ProxyServiceInformer) WatchServiceEvent() <-chan types.ServiceEvent {
	return i.ServiceChannel
}

func (i *ProxyServiceInformer) CloseChan() {
	close(i.ServiceChannel)
}

func (i *ProxyServiceInformer) ListServices() []object.Service {
	i.mtx.RLock()
	Services := make([]object.Service, len(i.ServiceCache))
	idx := 0
	for _, Service := range i.ServiceCache {
		Services[idx] = Service
		idx += 1
	}
	i.mtx.RUnlock()
	return Services
}

func (i *ProxyServiceInformer) InformService(newService object.Service, eType watchobj.EventType) error {
	i.mtx.Lock()
	oldService, exist := i.ServiceCache[newService.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.ServiceCache, newService.UID)
			i.ServiceChannel <- types.ServiceEvent{
				Type:    types.ServiceRemove,
				Service: newService,
			}
		} else {
			log.Printf("[INFO]: Service %s not exist, delete do nothing\n", newService.UID)
		}
	}

	if eType == watchobj.EVENT_PUT {
		// Just handle Service whose cluster ip is not empty
		if newService.Spec.ClusterIP == "" {
			log.Println("[INFO]: Service without cluster ip, just ignore")
			i.mtx.Unlock()
			return nil
		}

		// else cached and judge type
		i.ServiceCache[newService.UID] = newService
		if !exist {
			i.ServiceChannel <- types.ServiceEvent{
				Type:    types.ServiceCreate,
				Service: newService}
		} else {
			// compute Service change: IP / Label
			if object.ComputeServiceCriticalChange(&newService, &oldService) {
				log.Println("[INFO]: Service changed, Service ID is:", newService.UID)
				i.ServiceChannel <- types.ServiceEvent{
					Type:    types.ServiceUpdate,
					Service: newService}
			} else {
				log.Println("[INFO]: Service not changed, Service ID is:", newService.Name)
			}
		}
	}

	i.mtx.Unlock()
	return nil
}
