package controllermanager

import (
	"Cubernetes/pkg/controllermanager/controller/replicaset_controller"
)

type ControllerManager struct {
	RSController replicaset_controller.ReplicaSetController
	// other controller here
}

func NewControllerManager() (ControllerManager, error) {
	rsController, _ := replicaset_controller.NewReplicaSetController()
	return ControllerManager{
		RSController: rsController,
	}, nil
}

func (cm *ControllerManager) Run() {
	// watch ReplicaSet from API Server
	// then call Update / Remove function
}
