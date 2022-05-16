package cubeproxy

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"sync"
)

type Cubeproxy struct {
	//Runtime CubeproxyRuntime
	Runtime *proxyruntime.ProxyRuntime
	lock    sync.Mutex
}

func NewCubeProxy() *Cubeproxy {
	log.Printf("[INFO]: creating cubeproxy\n")
	runtime, err := proxyruntime.InitProxyRuntime()
	if err != nil {
		log.Printf("[Fatal]: Create cube proxy runtime error: %v", err.Error())
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

	// before watch service, add exist service to iptables
	pods, err := crudobj.GetPods()
	if err != nil {
		log.Fatalln("[Fatal]: get pods failed when init cubeproxy")
	}
	err = cp.Runtime.PodInformer.InitInformer(pods)
	if err != nil {
		log.Fatalln("[Fatal]: Init pod informer failed")
	}
	err = cp.Runtime.AddAllExistService()
	if err != nil {
		log.Fatalln("[Fatal]Add exist services failed")
	}

	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Println("[Error]: Error occurs when watching services")
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
				log.Panic("[Fatal]: Inform service failed")
				return
			}
		default:
			log.Panic("[Fatal]: Unsupported types in watching service.")
		}
	}

	log.Fatalln("[Fatal]: Unreachable here")
}

func (cp *Cubeproxy) syncService() {
	informEvent := cp.Runtime.ServiceInformer.WatchServiceEvent()

	for serviceEvent := range informEvent {
		log.Printf("[INFO]: [INFO]: Main loop working, types is %v,service id is %v", serviceEvent.Type, serviceEvent.Service.UID)
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
		log.Printf("[INFO]: Main loop working, type is %v, pod id is %v", podEvent.Type, &podEvent.Pod.UID)
		pod := podEvent.Pod
		eType := podEvent.Type
		cp.lock.Lock()

		switch eType {
		case types.PodCreate, types.PodRemove, types.PodUpdate:
			log.Printf("[INFO]: create pod %s\n", pod.UID)
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
		log.Printf("[INFO]: Main loop working, type is %v, DNS id is %v", podEvent.Type, &podEvent.DNS.UID)
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
			err := cp.Runtime.AddDNS(&dns)
			if err != nil {
				log.Fatalln("[Fatal]: error when modify DNS")
				return
			}
		}

		cp.lock.Unlock()
	}
}

func (cp *Cubeproxy) WatchPodsChange() error {
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Println("[Error]: Error occurs when watching pods")
		return err
	}
	defer cancel()

	for podEvent := range ch {
		switch podEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			err := cp.Runtime.PodInformer.InformPod(podEvent.Pod, podEvent.EType)
			if err != nil {
				log.Println("[Error]: Error when inform pod: ", podEvent.Pod.UID)
				return err
			}
		default:
			log.Panic("[Fatal]: Unsupported types in watch pod")
		}
	}

	log.Fatalln("[Fatal]: Unreachable here")
	return nil
}

func (cp *Cubeproxy) WatchDNSChange() error {
	ch, cancel, err := watchobj.WatchDnses()
	if err != nil {
		log.Println("[Error]: Error occurs when watching DNSes")
		return err
	}
	defer cancel()

	for dnsEvent := range ch {
		switch dnsEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			err := cp.Runtime.DNSInformer.InformDNS(dnsEvent.Dns, dnsEvent.EType)
			if err != nil {
				log.Println("[Error]: Error when inform DNS:", dnsEvent.Dns.UID)
				return err
			}
		default:
			log.Panic("[Fatal]: Unsupported types in watch dns")
		}
	}

	log.Fatalln("[Fatal]: Unreachable here")
	return nil
}
