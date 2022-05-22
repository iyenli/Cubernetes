package autoscaler_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

const (
	lowerReplicaSetNameSuffix     = "_REPLICA_SET"
	lowerReplocaSetParentUIDLabel = "cubernetes.replicaset.upper.uid"
)

func (asc *autoScalerController) handleAutoScalerCreate(as *object.AutoScaler) error {
	lowerRS := buildLowerReplicaSet(as)
	rs, err := crudobj.CreateReplicaSet(*lowerRS)
	if err != nil {
		log.Printf("fail to create ReplicaSet %s to API Server: %v\n", lowerRS.Name, err)
		return err
	}
	log.Printf("lower ReplicaSet of AutoScaler %s created\n", as.Name)

	as.Status = &object.AutoScalerStatus{
		LastScaleTime:     time.Now(),
		LastUpdateTime:    time.Now(),
		ReplicaSetUID:     rs.UID,
		DesiredReplicas:   int(rs.Spec.Replicas),
		ActualReplicas:    0,
		ActualUtilization: object.AverageUtilization{},
	}

	if _, err := crudobj.UpdateAutoScaler(*as); err != nil {
		log.Printf("fail to update autoscaler status to apiserver\n")
		return err
	}

	return nil
}

func (asc *autoScalerController) handleAutoScalerUpdate(as *object.AutoScaler) error {
	// Status update do nothing
	log.Printf("receive status update of AutoScaler %s\n", as.Name)
	return nil
}

func (asc *autoScalerController) handleAutoScalerRemove(as *object.AutoScaler) error {
	if err := crudobj.DeleteReplicaSet(as.Status.ReplicaSetUID); err != nil {
		log.Printf("fail to kill lower ReplicaSet of AutoScaler %s: %v\n", as.Name, err)
		return err
	}
	log.Printf("autoScaler %s removed\n", as.Name)
	return nil
}

func buildLowerReplicaSet(as *object.AutoScaler) *object.ReplicaSet {
	rs := &object.ReplicaSet{
		TypeMeta: object.TypeMeta{
			APIVersion: "wahtever/v1",
			Kind:       "ReplicaSet",
		},
		ObjectMeta: object.ObjectMeta{
			Name:        as.Name + lowerReplicaSetNameSuffix,
			Namespace:   as.Namespace,
			Labels:      as.Labels,
			Annotations: as.Annotations,
		},
		Spec: object.ReplicaSetSpec{
			// start with min replica
			Replicas: int32(as.Spec.MinReplicas),
			Selector: make(map[string]string),
			Template: as.Spec.Template,
		},
	}

	if rs.Labels == nil {
		rs.Labels = make(map[string]string)
	}
	// modify ReplicaSet LabelSelector & Pod template
	rs.Labels[lowerReplocaSetParentUIDLabel] = as.UID
	rs.Spec.Selector[lowerReplocaSetParentUIDLabel] = as.UID
	rs.Spec.Template.Labels[lowerReplocaSetParentUIDLabel] = as.UID

	return rs
}
