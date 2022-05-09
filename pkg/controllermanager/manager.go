package controllermanager

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/controller/replicaset_controller"
	"Cubernetes/pkg/controllermanager/informer"
	"Cubernetes/pkg/controllermanager/phase"
	"log"
)

type ControllerManager struct {
	RSController replicaset_controller.ReplicaSetController
	PodInformer  informer.PodInformer
	// other controller here
}

func NewControllerManager() ControllerManager {
	podInformer, _ := informer.NewPodInformer()
	rsController, _ := replicaset_controller.NewReplicaSetController(podInformer)
	return ControllerManager{
		RSController: rsController,
	}
}

func (cm *ControllerManager) Run() {
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		panic(err)
	}
	defer cancel()

	// watch ReplicaSet from API Server
	go cm.RSController.Run()

	for podEvent := range ch {
		if phase.NotHandle(podEvent.Pod.Status.Phase) {
			continue
		}
		switch podEvent.EType {
		case watchobj.EVENT_DELETE, watchobj.EVENT_PUT:
			cm.PodInformer.InformPod(&podEvent.Pod, podEvent.EType)
		default:
			log.Panic("Unsupported type in watch pod.")
		}
	}

	log.Fatalln("Unreachable here")
}
