package cubeproxy

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
)

type Cubeproxy struct {
	//Runtime CubeproxyRuntime
	Runtime *proxyruntime.ProxyRuntime
}

func (cp *Cubeproxy) syncLoop() {
	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Println("Error occurs when watching services")
		return
	}

	defer cancel()

	for serviceEvent := range ch {
		log.Printf("A service comes, type is %v, id is %v", serviceEvent.EType, serviceEvent.Service.UID)
		switch serviceEvent.EType {
		case watchobj.EVENT_PUT:
			err := cp.Runtime.AddService(&serviceEvent.Service)
			if err != nil {
				log.Printf("Add service error: %v", err.Error())
				return
			}
		case watchobj.EVENT_DELETE:
			err := cp.Runtime.DeleteService(&serviceEvent.Service)

			if err != nil {
				log.Printf("Delete service error: %v", err.Error())
				return
			}
		default:
			log.Panic("Unsupported type in watch service.")
		}
	}
}

func (cp *Cubeproxy) Run() {
	log.Println("Init IP Tables...")
	if cp.Runtime == nil {
		runtime, err := proxyruntime.InitIPTables()
		if err != nil {
			panic(err)
		}

		cp.Runtime = runtime
	}
	log.Println("Init IP Tables Success")

	defer func(runtime *proxyruntime.ProxyRuntime) {
		log.Printf("Release IP Tables...")
		err := runtime.ReleaseIPTables()
		if err != nil {
			log.Panicln("Error when release proxy Runtime")
		}
	}(cp.Runtime)

	cp.syncLoop()

	log.Fatalln("Unreachable here")
}
