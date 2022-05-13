package cubeproxy

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"sync"
)

type Cubeproxy struct {
	//Runtime CubeproxyRuntime
	Runtime *proxyruntime.ProxyRuntime

	lock sync.Mutex
}

func NewCubeProxy() *Cubeproxy {
	log.Printf("creating cubeproxy\n")
	runtime, err := proxyruntime.InitProxyRuntime()
	if err != nil {
		log.Printf("Create cube proxy runtime error: %v", err.Error())
	}

	cp := &Cubeproxy{
		Runtime: runtime,
		lock:    sync.Mutex{},
	}

	log.Println("Cubeproxy created")
	return cp
}

func (cp *Cubeproxy) Run() {
	if cp.Runtime == nil {
		log.Fatal("[Fatal]: Seg fault")
	}

	// TODO: Easy helper of cleaning iptables when exit unexpectedly
	defer func(runtime *proxyruntime.ProxyRuntime) {
		log.Printf("Release IP Tables...")
		err := runtime.ReleaseIPTables()
		if err != nil {
			log.Panicln("[Panic]: Error when release proxy Runtime")
		}
	}(cp.Runtime)

	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Println("Error occurs when watching services")
		return
	}
	defer cancel()

	// sync pod and service
	go cp.syncService()
	go cp.syncPod()

	// watch pod and service
	go func() {
		err := cp.WatchPodsChange()
		if err != nil {
			log.Fatalln("[Fatal]: watching pods in cubeproxy failed")
			return
		}
	}()

	for serviceEvent := range ch {
		log.Printf("A service comes, types is %v, id is %v", serviceEvent.EType, serviceEvent.Service.UID)
		switch serviceEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			err := cp.Runtime.ServiceInformer.InformService(serviceEvent.Service, serviceEvent.EType)
			if err != nil {
				log.Panic("Inform service failed")
				return
			}
		default:
			log.Panic("Unsupported types in watching service.")
		}
	}

	log.Fatalln("Unreachable here")
}

func (cp *Cubeproxy) syncService() {
	informEvent := cp.Runtime.ServiceInformer.WatchServiceEvent()

	for serviceEvent := range informEvent {
		log.Printf("Main loop working, types is %v,service id is %v", serviceEvent, serviceEvent.Service.UID)
		service := serviceEvent.Service
		eType := serviceEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.ServiceCreate:
			log.Printf("from serviceEvent: create service %s\n", service.UID)
			err := cp.Runtime.AddService(&service)
			if err != nil {
				log.Printf("Add service error: %v", err.Error())
				return
			}
		case types.ServiceUpdate:
			// critical update: simply delete and rebuild
			log.Printf("from serviceEvent: update service %s\n", service.UID)
			err := cp.Runtime.DeleteService(&service)
			if err != nil {
				log.Printf("Delete service error: %v", err.Error())
				return
			}

			err = cp.Runtime.AddService(&service)
			if err != nil {
				log.Printf("Add service error: %v", err.Error())
				return
			}

		case types.ServiceRemove:
			log.Printf("from serviceEvent: delete service %s\n", service.UID)
			err := cp.Runtime.DeleteService(&service)
			if err != nil {
				log.Printf("Delete service error: %v", err.Error())
				return
			}
		}

		cp.lock.Unlock()
	}
}

func (cp *Cubeproxy) syncPod() {
	informEvent := cp.Runtime.PodInformer.WatchPodEvent()

	for podEvent := range informEvent {
		log.Printf("Main loop working, types is %v,service id is %v", podEvent, &podEvent.Pod.UID)
		pod := podEvent.Pod
		eType := podEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.PodCreate, types.PodRemove, types.PodUpdate:
			log.Printf("from podEvent: create service %s\n", pod.UID)
			err := cp.Runtime.ModifyPod(&(pod))
			if err != nil {
				log.Fatalln("[Fatal]: error when modify pod")
				return
			}
		}

		cp.lock.Unlock()
	}
}

func (cp *Cubeproxy) WatchPodsChange() error {
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Println("Error occurs when watching pods")
		return err
	}
	defer cancel()

	for podEvent := range ch {
		log.Printf("A pod comes, types is %v, id is %v", podEvent.EType, podEvent.Pod.UID)
		switch podEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			err := cp.Runtime.PodInformer.InformPod(podEvent.Pod, podEvent.EType)
			if err != nil {
				log.Println("Error when inform pod: ", podEvent.Pod.UID)
				return err
			}
		default:
			log.Panic("Unsupported types in watch pod")
		}
	}

	log.Fatalln("Unreachable here")
	return nil
}
