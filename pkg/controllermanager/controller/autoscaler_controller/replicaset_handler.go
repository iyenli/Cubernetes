package autoscaler_controller

import (
	"Cubernetes/pkg/object"
	"log"
)

func (asc *autoScalerController) handleReplicaSetCreate(rs *object.ReplicaSet) error {
	// do nothing
	return nil
}

func (asc *autoScalerController) handleReplicaSetUpdate(rs *object.ReplicaSet) error {
	// do nothing
	return nil
}

// keep lower ReplicaSet alive in lifecycle of AutoScaler
// so we should create another ReplicaSet if old one was removed
func (asc *autoScalerController) handleReplicaSetRemove(rs *object.ReplicaSet) error {
	uid, ok := rs.Labels[lowerReplocaSetParentUIDLabel]
	if ok {
		as, found := asc.asInformer.GetAutoScaler(uid)
		if found {
			// create new ReplicaSet
			if err := asc.handleAutoScalerCreate(as); err != nil {
				log.Printf("fail to create new ReplicaSet for AutoScaler %s\n", as.Name)
				return err
			}
		}
	}
	return nil
}
