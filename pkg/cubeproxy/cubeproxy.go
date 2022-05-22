package cubeproxy

import (
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"sync"
)

type Cubeproxy struct {
	Runtime *proxyruntime.ProxyRuntime
	lock    sync.Mutex
}

func NewCubeProxy() *Cubeproxy {
	log.Println("[INFO]: creating cubeproxy")
	runtime, err := proxyruntime.InitProxyRuntime()
	if err != nil {
		log.Fatalf("[Fatal]: Create cube proxy runtime error: %v", err.Error())
	}

	cp := &Cubeproxy{
		Runtime: runtime,
		lock:    sync.Mutex{},
	}

	log.Println("[INFO]: Cubeproxy created")
	return cp
}

func (cp *Cubeproxy) Run() {
	if cp.Runtime == nil {
		log.Fatal("[Fatal]: Runtime not initialized before running")
	}

	defer func(runtime *proxyruntime.ProxyRuntime) {
		log.Printf("[INFO]: Release IP Tables...")
		err := runtime.ReleaseIPTables()
		if err != nil {
			log.Panicln("[Panic]: Error when release proxy Runtime")
		}
	}(cp.Runtime)

	err := cp.Runtime.AddAllExistService()
	if err != nil {
		log.Fatalln("[Fatal]Add exist services failed")
	}

	var wg sync.WaitGroup
	wg.Add(6)

	// sync pod and service and DNS
	go func() {
		defer wg.Done()
		cp.syncService()
	}()
	go func() {
		defer wg.Done()
		go cp.syncPod()
	}()
	go func() {
		defer wg.Done()
		go cp.syncDNS()
	}()

	// watch pod and service and DNS
	go func() {
		defer wg.Done()
		cp.Runtime.PodInformer.ListAndWatchPodsWithRetry()
	}()

	go func() {
		defer wg.Done()
		cp.Runtime.DNSInformer.ListAndWatchDNSWithRetry()
	}()

	go func() {
		defer wg.Done()
		cp.Runtime.ServiceInformer.ListAndWatchServicesWithRetry()
	}()

	wg.Wait()
	log.Fatalln("[Fatal]: Unreachable here")
}

func (cp *Cubeproxy) syncService() {
	informEvent := cp.Runtime.ServiceInformer.WatchServiceEvent()

	for serviceEvent := range informEvent {
		log.Printf("[INFO]: Main loop working, types is %v,service id is %v", serviceEvent.Type, serviceEvent.Service.UID)
		service := serviceEvent.Service
		eType := serviceEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.ServiceCreate:
			log.Printf("[INFO]: create service %s\n", service.UID)
			err := cp.Runtime.AddService(&service)
			if err != nil {
				log.Printf("[Error]: Add service error: %v", err.Error())
				return
			}
		case types.ServiceUpdate:
			// critical update: simply delete and rebuild
			log.Printf("[INFO]: update service %s\n", service.UID)
			err := cp.Runtime.DeleteService(&service)
			if err != nil {
				log.Printf("[Fatal]: Delete service error: %v", err.Error())
				return
			}

			err = cp.Runtime.AddService(&service)
			if err != nil {
				log.Printf("[Fatal]: Add service error: %v", err.Error())
				return
			}

		case types.ServiceRemove:
			log.Printf("[INFO]: delete service %s\n", service.UID)
			err := cp.Runtime.DeleteService(&service)
			if err != nil {
				log.Printf("[Fatal]: Delete service error: %v", err.Error())
				return
			}
		}

		cp.lock.Unlock()
	}
}

func (cp *Cubeproxy) syncPod() {
	informEvent := cp.Runtime.PodInformer.WatchPodEvent()

	for podEvent := range informEvent {
		log.Printf("[INFO]: Main loop working, type is %v, pod id is %v", podEvent.Type, podEvent.Pod.UID)
		pod := podEvent.Pod
		eType := podEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.PodCreate, types.PodRemove, types.PodUpdate:
			log.Printf("[INFO]: %v, podID is %s\n", eType, pod.UID)
			err := cp.Runtime.ModifyPod(&(pod))
			if err != nil {
				log.Fatalln("[Fatal]: error when modify pod")
				return
			}
		}

		cp.lock.Unlock()
	}
}

func (cp *Cubeproxy) syncDNS() {
	informEvent := cp.Runtime.DNSInformer.WatchDNSEvent()

	for podEvent := range informEvent {
		log.Printf("[INFO]: Main loop working, type is %v, DNS id is %v", podEvent.Type, podEvent.DNS.UID)
		dns := podEvent.DNS
		eType := podEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.DNSCreate:
			log.Printf("[INFO] DNS Created, DnsID %s\n", dns.UID)
			err := cp.Runtime.AddDNS(&dns)
			if err != nil {
				log.Fatalln("[Fatal]: error when create DNS")
				return
			}

		case types.DNSRemove:
			log.Printf("[INFO] DNS Removed, DnsID %s\n", dns.UID)
			err := cp.Runtime.DeleteDNS(&dns)
			if err != nil {
				log.Fatalln("[Fatal]: error when remove DNS")
				return
			}

		case types.DNSUpdate:
			log.Printf("[INFO] DNS Update, DnsID %s\n", dns.UID)
			err := cp.Runtime.ModifyDNS(&dns)
			if err != nil {
				log.Fatalln("[Fatal]: error when modify DNS")
				return
			}
		}

		cp.lock.Unlock()
	}
}
