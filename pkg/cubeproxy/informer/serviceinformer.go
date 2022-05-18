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

type ServiceInformer interface {
	ListAndWatchServicesWithRetry()
	WatchServiceEvent() <-chan types.ServiceEvent
	ListServices() []object.Service
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

func (i *ProxyServiceInformer) ListAndWatchServicesWithRetry() {
	defer close(i.ServiceChannel)
	for {
		i.tryListAndWatchServices()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (i *ProxyServiceInformer) tryListAndWatchServices() {
	if allServices, err := crudobj.GetServices(); err != nil {
		log.Printf("[INFO]: fail to get all services from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, service := range allServices {
			i.ServiceCache[service.UID] = service
		}
	}

	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Println("[Error]: Error occurs when watching services")
		return
	}
	defer cancel()

	for {
		select {
		case serviceEvent, ok := <-ch:
			{
				if !ok {
					log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
					return
				}
				log.Printf("A service comes, types is %v, id is %v", serviceEvent.EType, serviceEvent.Service.UID)
				switch serviceEvent.EType {
				case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
					err := i.informService(serviceEvent.Service, serviceEvent.EType)
					if err != nil {
						log.Panic("[Fatal]: Inform service failed")
						return
					}
				default:
					log.Panic("[Fatal]: Unsupported types in watching service.")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}

}

func (i *ProxyServiceInformer) WatchServiceEvent() <-chan types.ServiceEvent {
	return i.ServiceChannel
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

func (i *ProxyServiceInformer) informService(newService object.Service, eType watchobj.EventType) error {
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
