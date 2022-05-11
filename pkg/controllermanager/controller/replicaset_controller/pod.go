package replicaset_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

func (rsc *replicaSetController) handlePodCreate(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		rs.Status.RunningReplicas += 1
		rs.Status.PodUIDsRunning = append(rs.Status.PodUIDsRunning, pod.UID)
		if idx, found := rsc.podUIDAppearedIndex(pod.UID, rs.Status.PodUIDsToRun); found {
			rs.Status.PodUIDsToRun =
				append(rs.Status.PodUIDsToRun[:idx], rs.Status.PodUIDsToRun[idx+1:]...)
		} else {
			log.Printf("[FATAL] unexpected pod %s add to ReplicaSet %s when create\n", pod.Name, rs.Name)
		}

		rs.Status.LastUpdateTime = time.Now()
		if _, err := crudobj.UpdateReplicaSet(rs); err != nil {
			log.Printf("fail to update replicaset status to apiserver\n")
			return err
		}
	}
	return nil
}

func (rsc *replicaSetController) handlePodUpdate(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		// ReplicaSet update only needed when pod UID not expected
		if _, found := rsc.podUIDAppearedIndex(pod.UID, rs.Status.PodUIDsRunning); !found {
			rs.Status.RunningReplicas += 1
			rs.Status.PodUIDsRunning = append(rs.Status.PodUIDsRunning, pod.UID)
			log.Printf("[FATAL] unexpected pod %s add to ReplicaSet %s when update\n", pod.Name, rs.Name)

			if _, err := crudobj.UpdateReplicaSet(rs); err != nil {
				log.Printf("fail to update replicaset status to apiserver\n")
				return err
			}
		} else {
			log.Println("pod update do nothing")
		}
	}
	return nil
}

func (rsc *replicaSetController) handlePodKilled(pod *object.Pod) error {
	matched := rsc.rsInformer.GetMatchedReplicaSet(pod)
	for _, rs := range matched {
		if idx, found := rsc.podUIDAppearedIndex(pod.UID, rs.Status.PodUIDsToKill); found {
			rs.Status.PodUIDsToKill =
				append(rs.Status.PodUIDsToKill[:idx], rs.Status.PodUIDsToKill[idx+1:]...)
		} // else pod was killed by outside Cubernetes

		if idx, found := rsc.podUIDAppearedIndex(pod.UID, rs.Status.PodUIDsRunning); found {
			rs.Status.RunningReplicas -= 1
			rs.Status.PodUIDsRunning =
				append(rs.Status.PodUIDsRunning[:idx], rs.Status.PodUIDsRunning[idx+1:]...)
		} else {
			log.Printf("[FATAL] unexpected pod %s killed but not in running\n", pod.Name)
		}

		if _, err := crudobj.UpdateReplicaSet(rs); err != nil {
			log.Printf("fail to update replicaset status to apiserver\n")
			return err
		}
	}
	return nil
}

func (rsc *replicaSetController) getReplicaSetPods(rs *object.ReplicaSet) ([]object.Pod, error) {
	return rsc.podInformer.SelectPods(rs.Spec.Selector), nil
}

// index, found
func (rsc *replicaSetController) podUIDAppearedIndex(podUID string, list []string) (int, bool) {
	for idx, uid := range list {
		if uid == podUID {
			return idx, true
		}
	}
	return -1, false
}
