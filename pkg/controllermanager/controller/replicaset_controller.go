package controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
	"strings"

	"github.com/google/uuid"
)

type ReplicaSetController interface {
	UpdateReplicaSet(rs *object.ReplicaSet) error
	RemoveReplicaSet(UID string) error
}

const (
	podNameUUIDLen = 8
)

type replicaSetController struct {}

func NewReplicaSetController() (ReplicaSetController, error) {
	return &replicaSetController{}, nil
}

func (rsc *replicaSetController) UpdateReplicaSet(rs *object.ReplicaSet) error {

	currentPods, err := rsc.getReplicaSetPods(rs)
	if err != nil {
		log.Printf("fail to get pods by selector %v: %v\n", rs.Spec.Selector, err)
		return err
	}

	var toKeep, toKill, toCreate []*object.Pod
	for _, pod := range currentPods {
		if pod.Status.Phase != object.PodFailed && pod.Status.Phase != object.PodUnknown {
			toKeep = append(toKeep, &pod)
		} else {
			toKill = append(toKill, &pod)
		}
	}

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

func (rsc *replicaSetController) RemoveReplicaSet(UID string) error {
	// TODO: placeholder: get rs object from API Server
	var rs *object.ReplicaSet
	rs.Spec.Replicas = 0
	return rsc.UpdateReplicaSet(rs)
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

func (rsc *replicaSetController) getReplicaSetPods(rs *object.ReplicaSet) ([]object.Pod, error) {
	return crudobj.SelectPods(rs.Spec.Selector)
}
