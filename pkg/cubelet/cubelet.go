package cubelet

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/heartbeat"
	actruntime "Cubernetes/pkg/cubelet/actorruntime"
	"Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/cuberuntime"
	"Cubernetes/pkg/cubelet/gpuserver"
	"Cubernetes/pkg/cubelet/informer"
	informertypes "Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

type Cubelet struct {
	NodeID  string
	WeaveIP net.IP

	podInformer informer.PodInformer
	podRuntime  cuberuntime.CubeRuntime

	jobInformer informer.JobInformer
	jobRuntime  gpuserver.JobRuntime

	actorInformer informer.ActorInformer
	actorRuntime  actruntime.ActorRuntime

	bigLock sync.Mutex
}

func NewCubelet() *Cubelet {
	log.Printf("[INFO]: creating cubelet podRuntime manager\n")
	podRuntime, err := cuberuntime.NewCubeRuntimeManager()
	if err != nil {
		panic(err)
	}
	jobRuntime := gpuserver.NewJobRuntime()
	actorRuntime, _ := actruntime.NewActorRuntime()
	podInformer, _ := informer.NewPodInformer()
	jobInformer, _ := informer.NewJobInformer()
	actorInformer, _ := informer.NewActorInformer()

	log.Println("[INFO]: cubelet init ends")

	return &Cubelet{
		podInformer: podInformer,
		podRuntime:  podRuntime,

		jobInformer: jobInformer,
		jobRuntime:  jobRuntime,

		actorInformer: actorInformer,
		actorRuntime:  actorRuntime,
		bigLock:       sync.Mutex{},
	}
}

func (cl *Cubelet) InitCubelet(NodeUID string, ip net.IP) {
	log.Printf("Starting node, Node UID is %v, Node weave IP is %v", NodeUID, ip.String())
	cl.NodeID = NodeUID
	cl.WeaveIP = ip
	cl.podInformer.SetNodeUID(NodeUID)
	cl.jobInformer.SetNodeUID(NodeUID)
	cl.actorInformer.SetNodeUID(NodeUID)
}

func (cl *Cubelet) Run() {
	defer cl.podRuntime.Close()

	// push pod status to apiserver every 10 sec
	// simply using for loop to achieve block timer
	wg := sync.WaitGroup{}
	wg.Add(8)

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Second * 7)
			cl.updatePodsRoutine()
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Second * 14)
			cl.updateActorsRoutine()
		}
	}()

	// deal with pod event
	go func() {
		defer wg.Done()
		cl.syncPodLoop()
	}()

	go func() {
		defer wg.Done()
		cl.podInformer.ListAndWatchPodsWithRetry()
	}()

	go func() {
		defer wg.Done()
		cl.jobInformer.ListAndWatchJobsWithRetry()
	}()

	go func() {
		defer wg.Done()
		cl.actorInformer.ListAndWatchActorsWithRetry()
	}()

	go func() {
		defer wg.Done()
		cl.syncJobLoop()
	}()

	go func() {
		defer wg.Done()
		cl.syncActorLoop()
	}()

	wg.Wait()
	log.Fatalln("[Fatal]: sUnreachable here")
}

func (cl *Cubelet) syncPodLoop() {
	informEvent := cl.podInformer.WatchPodEvent()

	for podEvent := range informEvent {
		log.Printf("Main loop working, types is %v, pod id is %v", podEvent.Type, podEvent.Pod.UID)
		pod := podEvent.Pod
		eType := podEvent.Type
		cl.bigLock.Lock()

		switch eType {
		case informertypes.Create:
			log.Printf("[INFO]: podEvent coming: create pod %s\n", pod.UID)
			err := cl.podRuntime.SyncPod(&pod, &container.PodStatus{})
			if err != nil {
				log.Printf("fail to create pod %s: %v\n", pod.Name, err)
			}
		case informertypes.Update:
			log.Printf("[INFO]: podEvent coming: update pod %s\n", pod.UID)
			podStatus, err := cl.podRuntime.GetPodStatus(pod.UID)
			if err != nil {
				log.Printf("fail to get pod %s status: %v\n", pod.Name, err)
			}
			err = cl.podRuntime.SyncPod(&pod, podStatus)
			if err != nil {
				log.Printf("fail to update pod %s: %v\n", pod.Name, err)
			}
		case informertypes.Remove:
			err := cl.podRuntime.KillPod(pod.UID)
			if err != nil {
				log.Printf("fail to kill pod %s: %v\n", pod.Name, err)
			}
		}
		cl.bigLock.Unlock()
		// time.Sleep(time.Second * 2)
	}
}

func (cl *Cubelet) syncJobLoop() {
	informEvent := cl.jobInformer.WatchJobEvent()

	for jobEvent := range informEvent {
		log.Printf("[INFO]: main loop working, types is %v, job id is %v", jobEvent.Type, jobEvent.Job.UID)

		switch jobEvent.Type {
		case informertypes.Create:
			log.Printf("[INFO]: Event: create job %s\n", jobEvent.Job.UID)
			err := cl.jobRuntime.AddGPUJob(&jobEvent.Job)
			if err != nil {
				log.Printf("[Error]: fail to create job %s: %v\n", jobEvent, err)
			}
		default:
			log.Printf("[WARN]: Job only support adding now\n")
		}
	}
}

func (cl *Cubelet) syncActorLoop() {
	informEvent := cl.actorInformer.WatchActorEvent()

	for actorEvent := range informEvent {
		switch actorEvent.Type {
		case informertypes.Create:
			log.Printf("[INFO]: Event: create actor %s\n", actorEvent.Actor.UID)
			err := cl.actorRuntime.CreateActor(&actorEvent.Actor)
			if err != nil {
				log.Printf("[Error]: fail to create actor: %v\n", err)
			}
		case informertypes.Remove:
			log.Printf("[INFO]: Event: remove actor %s\n", actorEvent.Actor.UID)
			err := cl.actorRuntime.KillActor(actorEvent.Actor.UID)
			if err != nil {
				log.Printf("[Error]: fail to remove actor: %v\n", err)
			}
		case informertypes.Update:
			log.Printf("[INFO]: Event: update actor %s\n", actorEvent.Actor.UID)
			err := cl.actorRuntime.UpdateActionScript(&actorEvent.Actor)
			if err != nil {
				log.Printf("[Error]: fail to update actor: %v\n", err)
			}
		default:
			log.Printf("[WARN]: Actor only support Create & Kill now\n")
		}
	}
}

func (cl *Cubelet) updatePodsRoutine() {
	cl.bigLock.Lock()
	defer cl.bigLock.Unlock()

	if !heartbeat.CheckConn() {
		log.Printf("[FATAL] lost connection with apiserver: not update this time\n")
		return
	}

	// collect all pod in podCache
	pods := cl.podInformer.ListPods()

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
			podStatus, err := cl.podRuntime.InspectPod(&p)
			if err != nil {
				log.Printf("[Error]: fail to get pod status %s: %v\n", p.Name, err)
				podStatus = &object.PodStatus{Phase: object.PodUnknown}
			}

			podStatus.IP = ip
			podStatus.NodeUID = nodeUID
			log.Printf("[INFO]: updating pod status, ip is %v, status is %v",
				podStatus.IP.String(), podStatus.Phase)

			rp, err := crudobj.UpdatePodStatus(p.UID, *podStatus)
			if err != nil {
				status, _ := json.Marshal(*podStatus)
				log.Printf("[Error]: updating pod status, %v", string(status))
				log.Printf("[Error]: fail to push pod status %s: %v\n", p.UID, err)
				cl.podInformer.ForceRemove(p.UID)
			} else {
				log.Printf("[INFO]: push pod status %s: %s\n", rp.Name, podStatus.Phase)
			}
		}(pod, ip, nodeUID)
	}

	wg.Wait()
}

func (cl *Cubelet) updateActorsRoutine() {
	cl.bigLock.Lock()
	defer cl.bigLock.Unlock()

	if !heartbeat.CheckConn() {
		log.Printf("[FATAL] lost connection with apiserver: not update this time\n")
		return
	}

	actors := cl.actorInformer.ListActors()
	wg := sync.WaitGroup{}
	wg.Add(len(actors))

	for _, actor := range actors {
		if actor.Status.IP == nil {
			wg.Done()
			continue
		}

		go func(a object.Actor) {
			defer wg.Done()
			phase, err := cl.actorRuntime.InspectActor(a.UID)
			if err != nil {
				log.Printf("fail to inspect actor %s: %v", a.Name, err)
				return
			}

			a.Status.Phase = phase
			a.Status.LastUpdatedTime = time.Now()
			if _, err = crudobj.UpdateActor(a); err != nil {
				log.Printf("fail to update actor %s status: %v", a.Name, err)
			}
		}(actor)
	}

	wg.Wait()
}
