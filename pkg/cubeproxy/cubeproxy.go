package cubeproxy

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"log"
	"os"
)

type Cubeproxy struct {
	//Runtime CubeproxyRuntime
	Runtime *proxyruntime.ProxyRuntime
}

func (cp *Cubeproxy) syncLoop() {
	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Panic("Error occurs when watching services")
		os.Exit(0)
	}

	defer cancel()

	for serviceEvent := range ch {
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
	if cp.Runtime == nil {
		runtime, err := proxyruntime.InitIPTables()
		if err != nil {
			panic(err)
		}

		cp.Runtime = runtime
	}

	defer func(runtime *proxyruntime.ProxyRuntime) {
		err := runtime.ReleaseIPTables()
		if err != nil {
			log.Panicln("Error when release proxy Runtime")
		}
	}(cp.Runtime)

	cp.syncLoop()

	log.Fatalln("Unreachable here")
}
