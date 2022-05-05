package cubelet

import (
	watchobj "Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/container"
	cuberuntime "Cubernetes/pkg/cubelet/cuberuntime"
	"log"
	"os"
)

type Cubelet struct {
	NodeID  string
	runtime cuberuntime.CubeRuntime
}

func NewCubelet() *Cubelet {
	log.Printf("creating cubelet runtime manager\n")
	runtime, err := cuberuntime.NewCubeRuntimeManager()
	if err != nil {
		panic(err)
	}

	return &Cubelet{runtime: runtime}
}

func (cl *Cubelet) Run() {
	defer cl.runtime.Close()
	cl.syncLoop()
	log.Fatalln("Unreachable here")
}

func (cl *Cubelet) syncLoop() {
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Panic("Error occurs when watching pods")
		os.Exit(0)
	}

	defer cancel()

	for podEvent := range ch {
		switch podEvent.EType {
		case watchobj.EVENT_PUT:
			err := cl.runtime.SyncPod(&podEvent.Pod, &container.PodStatus{})
			if err != nil {
				log.Printf("error when sync pod: %v", err)
			}
		case watchobj.EVENT_DELETE:
			err := cl.runtime.KillPod(podEvent.Pod.UID)
			if err != nil {
				log.Printf("error when delete pod: %v", err)
			}
		default:
			log.Panic("Unsupported type in watch pod.")
		}
	}
}

func (cl *Cubelet) updatePodPeriod() {

}
