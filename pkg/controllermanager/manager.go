package controllermanager

import "Cubernetes/pkg/controllermanager/controller"

type ControllerManager struct {
	RSController controller.ReplicaSetController
	// other controller here
}

func NewControllerManager() (ControllerManager, error) {
	rsController, _ := controller.NewReplicaSetController()
	return ControllerManager{
		RSController: rsController,
	}, nil
}

func (cm *ControllerManager) Run() {
	// watch ReplicaSet from API Server
	// then call Update / Remove function
}