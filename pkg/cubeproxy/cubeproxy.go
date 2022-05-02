package cubeproxy

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"log"
	"os"
)

type Cubeproxy struct {
	//runtime CubeproxyRuntime
}

func syncLoop() {
	ch, cancel, err := watchobj.WatchServices()
	if err != nil {
		log.Panic("Error occurs when watching services")
		os.Exit(0)
	}

	defer cancel()

	for serviceEvent := range ch {
		switch serviceEvent.EType {
		case watchobj.EVENT_PUT:
			log.Println("TODO")
		case watchobj.EVENT_DELETE:
			log.Println("TODO")

		default:
			log.Panic("Unsupported type in watch service.")
		}
	}
}

func Run() {

	syncLoop()
}
