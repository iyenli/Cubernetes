package replicaset_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/health"
	"Cubernetes/pkg/controllermanager/informer"
	"Cubernetes/pkg/controllermanager/phase"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/controllermanager/utils"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

const (
	replicaSetPodsUpdateWaitTime = time.Second * 30
)

type ReplicaSetController interface {
	Run()
}

type replicaSetController struct {
	podInformer informer.PodInformer
	rsInformer  informer.ReplicaSetInformer
	biglock     sync.Mutex
	wg          sync.WaitGroup
}

func NewReplicaSetController(
	podInformer informer.PodInformer,
	rsInformer informer.ReplicaSetInformer,
	wg sync.WaitGroup) (ReplicaSetController, error) {
	wg.Add(1)
	return &replicaSetController{
		podInformer: podInformer,
		rsInformer:  rsInformer,
		biglock:     sync.Mutex{},
		wg:          wg,
	}, nil
}

func (rsc *replicaSetController) Run() {

	go func() {
		for {
			time.Sleep(time.Second * 10)
			rsc.updateReplicaSetsRoutine()
		}
	}()

	rsc.syncLoop()
}

func (rsc *replicaSetController) syncLoop() {
	podEventChan := rsc.podInformer.WatchPodEvent()
	defer rsc.podInformer.CloseChan(podEventChan)

	rsEventChan := rsc.rsInformer.WatchRSEvent()
	defer rsc.rsInformer.CloseChan(rsEventChan)

	rsc.wg.Done()

	for {
		select {
		case podEvent := <-podEventChan:
			rsc.biglock.Lock()
			pod := podEvent.Pod
			switch podEvent.Type {
			case types.PodCreate:
				log.Printf("handle Pod %s create\n", podEvent.Pod.Name)
				rsc.handlePodCreate(&pod)
			case types.PodUpdate:
				log.Printf("handle Pod %s update\n", podEvent.Pod.Name)
				rsc.handlePodUpdate(&pod)
			case types.PodKilled:
				log.Printf("handle Pod %s killed\n", podEvent.Pod.Name)
				rsc.handlePodKilled(&pod)
			default:
				log.Fatal("[FATAL] Unknown podInformer event types: " + podEvent.Type)
			}
			rsc.biglock.Unlock()
		case rsEvent := <-rsEventChan:
			rsc.biglock.Lock()
			replicaSet := rsEvent.ReplicaSet
			switch rsEvent.Type {
			case types.RsCreate:
				log.Printf("handle ReplicaSet %s create\n", rsEvent.ReplicaSet.Name)
				rsc.handleReplicaSetCreate(&replicaSet)
			case types.RsUpdate:
				log.Printf("handle ReplicaSet %s update\n", rsEvent.ReplicaSet.Name)
				rsc.handleReplicaSetUpdate(&replicaSet)
			case types.RsRemove:
				log.Printf("handle ReplicaSet %s remove\n", rsEvent.ReplicaSet.Name)
				rsc.handleReplicaSetRemove(&replicaSet)
			default:
				log.Fatal("[FATAL] Unknown rsInformer event types: " + rsEvent.Type)
			}
			rsc.biglock.Unlock()
		default:
			time.Sleep(time.Second * 4)
		}
	}

}

func (rsc *replicaSetController) updateReplicaSetsRoutine() {
	rsc.biglock.Lock()
	defer rsc.biglock.Unlock()

	if !health.CheckApiServerHealth() {
		log.Printf("[FATAL] lost connection with apiserver: not update this time\n")
		return
	}

	replicaSets := rsc.rsInformer.ListReplicaSets()

	wg := sync.WaitGroup{}
	wg.Add(len(replicaSets))
	for _, replicaSet := range replicaSets {
		go func(rs object.ReplicaSet) {
			defer wg.Done()
			// rs.Status will be nil if a ReplicaSet is just created by cubectl
			// but not yet handled by handleReplicaSetCreate
			if rs.Status != nil {
				rsc.checkAndUpdateReplicaSetStatus(&rs)
			}
		}(replicaSet)
	}
	wg.Wait()
}

func (rsc *replicaSetController) checkAndUpdateReplicaSetStatus(rs *object.ReplicaSet) error {

	if time.Since(rs.Status.LastUpdateTime) < replicaSetPodsUpdateWaitTime {
		return nil
	}

	currentPods, err := rsc.getReplicaSetPods(rs)
	if err != nil {
		log.Printf("fail to get pods by selector %v: %v\n", rs.Spec.Selector, err)
		return err
	}

	runnings := make([]string, 0)
	bads := make([]string, 0)
	for _, pod := range currentPods {
		if phase.Running(pod.Status.Phase) {
			runnings = append(runnings, pod.UID)
		} else if phase.Bad(pod.Status.Phase) {
			bads = append(bads, pod.UID)
		}
	}

	log.Printf("check %s status:\n", rs.Name)
	log.Println("toRun:   ", rs.Status.PodUIDsToRun)
	log.Println("toKill:  ", rs.Status.PodUIDsToKill)
	log.Println("running: ", runnings)
	log.Println("bad:     ", bads)

	// timeout => create new pods, ignore old
	toCreate := int(rs.Spec.Replicas) - len(runnings)
	log.Printf("%d pod(s) running, %d expected\n", len(runnings), rs.Spec.Replicas)
	podsToRun := make([]string, 0)
	// will do nothing if toCreate <= 0
	for idx := 0; idx < toCreate; idx += 1 {
		newPod := rsc.buildNewAPIPod(&rs.Spec.Template, rs.Name)
		if pod, err := crudobj.CreatePod(*newPod); err != nil {
			log.Printf("fail to create pod %s to API Server: %v\n", newPod.Name, err)
		} else {
			log.Printf("ReplicaSet %s add pod: %s (%s)\n", rs.Name, pod.Name, pod.UID)
			podsToRun = append(podsToRun, pod.UID)
		}
	}

	// timeout => kill toKill again, and try to kill old-fail-to-create
	podsToKill := append(rs.Status.PodUIDsToKill, rs.Status.PodUIDsToRun...)
	// kill redundant pod if needed
	if toCreate < 0 {
		podsToKill = append(podsToKill, runnings[:-toCreate]...)
	}
	// also kill bad pods found in update, then remove duplication
	podsToKill = utils.RemoveDuplication(append(podsToKill, bads...))
	for _, uid := range podsToKill {
		if err := crudobj.DeletePod(uid); err != nil {
			log.Printf("fail to delete pod %s from API Server: %v\n", uid, err)
		} else {
			log.Printf("ReplicaSet %s remove pod from API Server: %s\n", rs.Name, uid)
		}
	}

	rs.Status = &object.ReplicaSetStatus{
		RunningReplicas: int32(len(runnings)),
		PodUIDsToRun:    podsToRun,
		PodUIDsToKill:   podsToKill,
		PodUIDsRunning:  runnings,
		LastUpdateTime:  time.Now(),
	}

	if _, err := crudobj.UpdateReplicaSet(*rs); err != nil {
		log.Printf("fail to update replicaset status to apiserver\n")
		return err
	}

	return nil
}
