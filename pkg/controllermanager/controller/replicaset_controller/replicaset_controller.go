package replicaset_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/informer"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReplicaSetController interface {
	Run()
}

const (
	podNameUUIDLen = 8
)

type replicaSetController struct {
	podInformer informer.PodInformer
	rsInformer  ReplicaSetInformer
}

func NewReplicaSetController(podInformer informer.PodInformer) (ReplicaSetController, error) {
	rsInformer, _ := NewReplicaSetInformer()
	return &replicaSetController{
		podInformer: podInformer,
		rsInformer:  rsInformer,
	}, nil
}

func (rsc *replicaSetController) Run() {
	ch, cancel, err := watchobj.WatchReplicaSets()
	if err != nil {
		log.Printf("fail to watch ReplicaSets from apiserver: %v\n", err)
		return
	}
	defer cancel()

	go rsc.syncLoop()

	for rsEvent := range ch {
		switch rsEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			rsc.rsInformer.InformReplicaSet(&rsEvent.ReplicaSet, rsEvent.EType)
		default:
			log.Fatal("[FATAL] Unknown event type: " + rsEvent.EType)
		}
	}
}

func (rsc *replicaSetController) syncLoop() {
	podEventChan := rsc.podInformer.WatchPodEvent()
	defer rsc.podInformer.CloseChan(podEventChan)

	rsEventChan := rsc.rsInformer.WatchRSEvent()
	defer rsc.rsInformer.CloseChan()

	select {
	case podEvent := <-podEventChan:
		switch podEvent.Type {
		case types.PodCreate:
			rsc.handlePodCreate(podEvent.Pod)
		case types.PodUpdate:
			rsc.handlePodUpdate(podEvent.Pod)
		case types.PodKilled:
			rsc.handlePodKilled(podEvent.Pod)
		default:
			log.Fatal("[FATAL] Unknown podInformer event type: " + podEvent.Type)
		}
	case rsEvent := <-rsEventChan:
		switch rsEvent.Type {
		case rsCreate, rsUpdate:
			rsc.updateReplicaSet(rsEvent.ReplicaSet)
		case rsRemove:
			rsc.removeReplicaSet(rsEvent.ReplicaSet)
		default:
			log.Fatal("[FATAL] Unknown rsInformer event type: " + rsEvent.Type)
		}
	default:
		time.Sleep(time.Second * 3)
	}

}

func (rsc *replicaSetController) updateReplicaSet(rs *object.ReplicaSet) error {

	currentPods, err := rsc.getReplicaSetPods(rs)
	if err != nil {
		log.Printf("fail to get pods by selector %v: %v\n", rs.Spec.Selector, err)
		return err
	}

	var toKeep, toKill, toCreate []*object.Pod
	for _, pod := range currentPods {
		if pod.Status.Phase != object.PodFailed && pod.Status.Phase != object.PodUnknown {
			toKeep = append(toKeep, pod)
		} else {
			toKill = append(toKill, pod)
		}
	}

	// update ReplicaSet status to apiserver
	rs.Status.RunningReplicas = int32(len(toKeep))
	rs.Status.PodUIDs = make([]string, len(toKeep))
	for idx, pod := range toKeep {
		rs.Status.PodUIDs[idx] = pod.UID
	}
	crudobj.UpdateReplicaSet(*rs)

	if len(toKeep) > int(rs.Spec.Replicas) {
		// kill redundant replica
		keepCount := rs.Spec.Replicas
		toKill = append(toKill, toKeep[keepCount:]...)
	} else {
		// add new API pod to API Server
		createCount := int(rs.Spec.Replicas) - len(toKeep)
		for i := 0; i < createCount; i += 1 {
			toCreate = append(toKeep, rsc.buildNewAPIPod(&rs.Spec.Template, rs.Name))
		}
	}

	for _, kill := range toKill {
		if err = crudobj.DeletePod(kill.UID); err != nil {
			log.Printf("fail to delete pod %s from API Server: %v\n", kill.UID, err)
		}
	}

	for _, create := range toCreate {
		if pod, err := crudobj.CreatePod(*create); err != nil {
			log.Printf("fail to create pod %s to API Server: %v\n", create.Name, err)
		} else {
			log.Printf("ReplicaSet %s add new pod to API Server: %v\n", rs.Name, pod)
		}
	}

	return err
}

func (rsc *replicaSetController) removeReplicaSet(rs *object.ReplicaSet) error {
	// just set its replica to 0
	rs.Spec.Replicas = 0
	return rsc.updateReplicaSet(rs)
}

// TODO: using the same handler now

func (rsc *replicaSetController) handlePodCreate(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		rsc.updateReplicaSet(rs)
	}
	return nil
}

func (rsc *replicaSetController) handlePodUpdate(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		rsc.updateReplicaSet(rs)
	}
	return nil
}

func (rsc *replicaSetController) handlePodKilled(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		rsc.updateReplicaSet(rs)
	}
	return nil
}

func (rsc *replicaSetController) buildNewAPIPod(template *object.PodTemplate, scName string) *object.Pod {

	pod := &object.Pod{
		TypeMeta: object.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: template.ObjectMeta,
		Spec:       template.Spec,
	}

	pod.Name = rsc.buildTemplatePodName(scName)

	return pod
}

func (rsc *replicaSetController) buildTemplatePodName(scName string) string {
	return strings.Join([]string{scName, uuid.New().String()[:podNameUUIDLen]}, "_")
}

func (rsc *replicaSetController) getReplicaSetPods(rs *object.ReplicaSet) ([]*object.Pod, error) {
	return rsc.podInformer.SelectPods(rs.Spec.Selector), nil
}
