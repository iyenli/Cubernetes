package cubelet

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/apiserver/crudobj"
	watchobj "Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/container"
	cuberuntime "Cubernetes/pkg/cubelet/cuberuntime"
	"Cubernetes/pkg/cubelet/informer"
	informertypes "Cubernetes/pkg/cubelet/informer/types"
	"log"
	"sync"
	"time"
)

type Cubelet struct {
	NodeID   string
	informer informer.PodInformer
	runtime  cuberuntime.CubeRuntime
}

func NewCubelet() *Cubelet {
	log.Printf("creating cubelet runtime manager\n")
	runtime, err := cuberuntime.NewCubeRuntimeManager()
	if err != nil {
		panic(err)
	}

	podInformer, _ := informer.NewPodInformer()
	log.Println("cubelet init ends")
	return &Cubelet{
		informer: podInformer,
		runtime:  runtime,
	}
}

func (cl *Cubelet) InitCubelet(NodeUID string) {
	log.Println("Starting node, Node UID is ", NodeUID)
	cubeconfig.NodeUID = NodeUID
}

func (cl *Cubelet) Run() {
	defer cl.runtime.Close()
	defer cl.informer.CloseChan()

	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Panic("Error occurs when watching pods")
	}

	defer cancel()

	// push pod status to apiserver every 10 sec
	// simply using for loop to achieve block timer
	go func() {
		for {
			time.Sleep(time.Second * 10)
			cl.updatePodsPeriod()
		}
	}()

	// deal with pod event
	go cl.syncLoop()

	for podEvent := range ch {
		if podEvent.Pod.Status == nil {
			log.Println("[INFO] Pod Catch, but status is nil so Cubelet don't handle")
			continue
		}
		if podEvent.Pod.Status.PodUID == cubeconfig.NodeUID {
			log.Println("[INFO] my pod Catch, type is ", podEvent.EType)

			switch podEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				err := cl.informer.InformPod(&podEvent.Pod, podEvent.EType)
				if err != nil {
					return
				}
			default:
				log.Panic("Unsupported type in watch pod.")
			}
		} else {
			log.Printf("[INFO] Pod Catch, but not my pod, pod UUID = %v, my UUID = %v", podEvent.Pod.Status.PodUID, cubeconfig.NodeUID)
		}
	}

	log.Fatalln("Unreachable here")
}

func (cl *Cubelet) syncLoop() {
	informEvent := cl.informer.PodEvent()

	for podEvent := range informEvent {
		log.Printf("Main loop working, type is %v, pod id is %v", podEvent.Type, podEvent.Pod.UID)
		switch podEvent.Type {
		case informertypes.PodCreate:
			err := cl.runtime.SyncPod(podEvent.Pod, &container.PodStatus{})
			if err != nil {
				log.Printf("fail to create pod %s: %v\n", podEvent.Pod.Name, err)
			}
		case informertypes.PodUpdate:
			podStatus, err := cl.runtime.GetPodStatus(podEvent.Pod.UID)
			if err != nil {
				log.Printf("fail to get pod %s status: %v\n", podEvent.Pod.Name, err)
			}
			err = cl.runtime.SyncPod(podEvent.Pod, podStatus)
			if err != nil {
				log.Printf("fail to update pod %s: %v\n", podEvent.Pod.Name, err)
			}
		case informertypes.PodRemove:
			err := cl.runtime.KillPod(podEvent.Pod.UID)
			if err != nil {
				log.Printf("fail to kill pod %s: %v\n", podEvent.Pod.Name, err)
			}
		}
		time.Sleep(time.Second * 2)
	}
}

func (cl *Cubelet) updatePodsPeriod() {

	// collect all pod by its sandbox
	uids, err := cl.runtime.ListPodsUID()
	if err != nil {
		log.Printf("fail to list uid of all pods\n")
	}

	// parallelly push all pod status to apiserver
	wg := sync.WaitGroup{}
	wg.Add(len(uids))

	for _, podUID := range uids {
		go func(uid string) {
			defer wg.Done()
			podStatus, err := cl.runtime.InspectPod(uid)
			if err != nil {
				log.Printf("fail to get pod status %s: %v\n", uid, err)
			}
			pod, err := crudobj.UpdatePodStatus(uid, *podStatus)
			if err != nil {
				log.Printf("fail to push pod status %s: %v\n", uid, err)
			} else {
				log.Printf("push pod status %s: %s\n", pod.Name, podStatus.Phase)
			}
		}(podUID)
	}

	wg.Wait()
}
