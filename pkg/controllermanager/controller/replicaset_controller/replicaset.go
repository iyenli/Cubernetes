package replicaset_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	podNameUUIDLen = 8
)

func (rsc *replicaSetController) handleReplicaSetCreate(rs *object.ReplicaSet) error {
	toCreate := rs.Spec.Replicas
	podsToRun := make([]string, 0)
	var err error

	for idx := 0; idx < int(toCreate); idx += 1 {
		newPod := rsc.buildNewAPIPod(&rs.Spec.Template, rs.Name)
		if pod, err := crudobj.CreatePod(*newPod); err != nil {
			log.Printf("fail to create pod %s to API Server: %v\n", newPod.Name, err)
		} else {
			log.Printf("ReplicaSet %s add pod: %s(%s)\n", rs.Name, pod.Name, pod.UID)
			podsToRun = append(podsToRun, pod.UID)
		}
	}

	rs.Status = &object.ReplicaSetStatus{
		RunningReplicas: 0,
		PodUIDsToRun:    podsToRun,
		LastUpdateTime:  time.Now(),
	}

	if _, err = crudobj.UpdateReplicaSet(*rs); err != nil {
		log.Printf("fail to update replicaset status to apiserver\n")
		return err
	}

	return nil
}

func (rsc *replicaSetController) handleReplicaSetUpdate(rs *object.ReplicaSet) error {
	// only handle replicas number update:
	// Template Spec update will handled by remove + create

	if int(rs.Spec.Replicas) > len(rs.Status.PodUIDsRunning) {
		toCreate := int(rs.Spec.Replicas) - len(rs.Status.PodUIDsRunning)
		for idx := 0; idx < toCreate; idx += 1 {
			newPod := rsc.buildNewAPIPod(&rs.Spec.Template, rs.Name)
			if pod, err := crudobj.CreatePod(*newPod); err != nil {
				log.Printf("fail to create pod %s to API Server: %v\n", newPod.Name, err)
			} else {
				log.Printf("ReplicaSet %s add pod: %s (%s)\n", rs.Name, pod.Name, pod.UID)
				rs.Status.PodUIDsToRun = append(rs.Status.PodUIDsToRun, pod.UID)
			}
		}
	} else {
		toKill := len(rs.Status.PodUIDsRunning) - int(rs.Spec.Replicas)
		rs.Status.PodUIDsToKill = append(rs.Status.PodUIDsToKill, rs.Status.PodUIDsRunning[:toKill]...)

		for _, uid := range rs.Status.PodUIDsRunning[:toKill] {
			if err := crudobj.DeletePod(uid); err != nil {
				log.Printf("fail to delete pod %s from API Server: %v\n", uid, err)
			} else {
				log.Printf("ReplicaSet %s remove pod from API Server: %s\n", rs.Name, uid)
			}
		}
	}

	rs.Status.LastUpdateTime = time.Now()
	if _, err := crudobj.UpdateReplicaSet(*rs); err != nil {
		log.Printf("fail to update replicaset status to apiserver\n")
		return err
	}

	return nil
}

func (rsc *replicaSetController) handleReplicaSetRemove(toKill *object.ReplicaSet) error {
	toKill.Spec.Replicas = 0
	return rsc.handleReplicaSetUpdate(toKill)
}

func (rsc *replicaSetController) buildNewAPIPod(template *object.PodTemplate, scName string) *object.Pod {

	pod := &object.Pod{
		TypeMeta: object.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: template.ObjectMeta,
		Spec:       template.Spec,
		Status: &object.PodStatus{
			Phase: object.PodCreated,
		},
	}

	pod.Name = rsc.buildTemplatePodName(scName)

	return pod
}

func (rsc *replicaSetController) buildTemplatePodName(scName string) string {
	return strings.Join([]string{scName, uuid.New().String()[:podNameUUIDLen]}, "_")
}
