package cubelet

import (
	"Cubernetes/pkg/apiserver/crudobj"
	watchobj "Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/container"
	cuberuntime "Cubernetes/pkg/cubelet/cuberuntime"
	"Cubernetes/pkg/cubelet/informer"
	informertypes "Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

type Cubelet struct {
	NodeID   string
	informer informer.PodInformer
	runtime  cuberuntime.CubeRuntime
	biglock  sync.Mutex
}

func NewCubelet() *Cubelet {
	log.Printf("creating cubelet runtime manager\n")
	runtime, err := cuberuntime.NewCubeRuntimeManager()
	if err != nil {
		panic(err)
	}

	informer, _ := informer.NewPodInformer()

	return &Cubelet{
		informer: informer,
		runtime:  runtime,
		biglock:  sync.Mutex{},
	}
}

func (cl *Cubelet) Run() {
	defer cl.runtime.Close()
	defer cl.informer.CloseChan()

	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Panic("Error occurs when watching pods")
		panic(err)
	}

	defer cancel()

	// push pod status to apiserver every 10 sec
	// simply using for loop to achieve block timer
	go func() {
		for {
			time.Sleep(time.Second * 7)
			cl.updatePodsRoutine()
		}
	}()

	// deal with pod event
	go cl.syncLoop()

	for podEvent := range ch {
		switch podEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			cl.informer.InformPod(podEvent.Pod, podEvent.EType)
		default:
			log.Panic("Unsupported type in watch pod.")
		}
	}

	log.Fatalln("Unreachable here")
}

func (cl *Cubelet) syncLoop() {
	informEvent := cl.informer.WatchPodEvent()

	for podEvent := range informEvent {
		pod := podEvent.Pod
		eType := podEvent.Type
		cl.biglock.Lock()
		switch eType {
		case informertypes.PodCreate:
			log.Printf("from podEvent: create pod %s\n", pod.UID)
			err := cl.runtime.SyncPod(&pod, &container.PodStatus{})
			if err != nil {
				log.Printf("fail to create pod %s: %v\n", pod.Name, err)
			}
		case informertypes.PodUpdate:
			log.Printf("from podEvent: update pod %s\n", pod.UID)
			podStatus, err := cl.runtime.GetPodStatus(pod.UID)
			if err != nil {
				log.Printf("fail to get pod %s status: %v\n", pod.Name, err)
			}
			err = cl.runtime.SyncPod(&pod, podStatus)
			if err != nil {
				log.Printf("fail to update pod %s: %v\n", pod.Name, err)
			}
		case informertypes.PodRemove:
			err := cl.runtime.KillPod(pod.UID)
			if err != nil {
				log.Printf("fail to kill pod %s: %v\n", pod.Name, err)
			}
		}
		cl.biglock.Unlock()
		// time.Sleep(time.Second * 2)
	}
}

func (cl *Cubelet) updatePodsRoutine() {
	cl.biglock.Lock()
	defer cl.biglock.Unlock()

	// collect all pod in podCache
	pods := cl.informer.ListPods()

	// parallelly push all pod status to apiserver
	wg := sync.WaitGroup{}
	wg.Add(len(pods))

	for _, pod := range pods {
		go func(p object.Pod) {
			defer wg.Done()
			podStatus, err := cl.runtime.InspectPod(&p)
			if err != nil {
				log.Printf("fail to get pod status %s: %v\n", p.Name, err)
				podStatus = &object.PodStatus{Phase: object.PodUnknown}
			}
			rp, err := crudobj.UpdatePodStatus(p.UID, *podStatus)
			if err != nil {
				log.Printf("fail to push pod status %s: %v\n", p.UID, err)
			} else {
				log.Printf("push pod status %s: %s\n", rp.Name, podStatus.Phase)
			}
		}(pod)
	}

	wg.Wait()
}
