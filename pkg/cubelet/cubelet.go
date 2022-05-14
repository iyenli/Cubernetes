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
	"net"
	"sync"
	"time"
)

type Cubelet struct {
	NodeID string
	// in weave container, this ip is host ip
	WeaveIP  net.IP
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

	podInformer, _ := informer.NewPodInformer()
	log.Println("cubelet init ends")

	return &Cubelet{
		informer: podInformer,
		runtime:  runtime,
		biglock:  sync.Mutex{},
	}
}

func (cl *Cubelet) InitCubelet(NodeUID string, ip net.IP) {
	log.Printf("Starting node, Node UID is %v, Node weave IP is %v", NodeUID, ip.String())
	cl.NodeID = NodeUID
	cl.WeaveIP = ip
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
			time.Sleep(time.Second * 7)
			cl.updatePodsRoutine()
		}
	}()

	// deal with pod event
	go cl.syncLoop()

	for podEvent := range ch {
		if podEvent.Pod.Status == nil && podEvent.EType != watchobj.EVENT_DELETE {
			log.Println("[INFO] Pod caught, but status is nil so Cubelet doesn't handle it")
			continue
		}
		if podEvent.EType == watchobj.EVENT_DELETE || podEvent.Pod.Status.NodeUID == cl.NodeID {
			log.Println("[INFO]: my pod caught, types is", podEvent.EType)
			switch podEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				err := cl.informer.InformPod(podEvent.Pod, podEvent.EType)
				if err != nil {
					return
				}
			default:
				log.Panic("Unsupported types in watch pod.")
			}
		} else {
			log.Printf("[INFO]: pod caught, but not my pod, pod UUID = %v, my UUID = %v",
				podEvent.Pod.Status.NodeUID, cl.NodeID)
		}
	}

	log.Fatalln("Unreachable here")
}

func (cl *Cubelet) syncLoop() {
	informEvent := cl.informer.WatchPodEvent()

	for podEvent := range informEvent {
		log.Printf("Main loop working, types is %v, pod id is %v", podEvent.Type, podEvent.Pod.UID)
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

	// parallel push all pod status to apiserver
	wg := sync.WaitGroup{}
	wg.Add(len(pods))

	for _, pod := range pods {
		ip := pod.Status.IP
		nodeUID := pod.Status.NodeUID
		log.Printf("[INFO]: Ready to update pod, ip is %v, nodeID is %v",
			ip.String(), nodeUID)

		if ip == nil {
			log.Printf("[INFO]: Not updating pod(%v) status before IP allocating\n", pod.UID)
			wg.Done()
			continue
		}

		go func(p object.Pod, ip net.IP, uid string) {
			defer wg.Done()
			podStatus, err := cl.runtime.InspectPod(&p)
			if err != nil {
				log.Printf("fail to get pod status %s: %v\n", p.Name, err)
				podStatus = &object.PodStatus{Phase: object.PodUnknown}
			}

			podStatus.IP = ip
			podStatus.NodeUID = nodeUID
			log.Printf("[INFO]: updating pod status, ip is %v, status is %v",
				podStatus.IP.String(), podStatus.Phase)

			rp, err := crudobj.UpdatePodStatus(p.UID, *podStatus)
			if err != nil {
				log.Printf("fail to push pod status %s: %v\n", p.UID, err)
			} else {
				log.Printf("[INFO]: push pod status %s: %s\n", rp.Name, podStatus.Phase)
			}
		}(pod, ip, nodeUID)
	}

	wg.Wait()
}
