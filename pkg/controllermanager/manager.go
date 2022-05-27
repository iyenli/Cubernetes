package controllermanager

import (
	"Cubernetes/pkg/controllermanager/controller/autoscaler_controller"
	"Cubernetes/pkg/controllermanager/controller/replicaset_controller"
	"Cubernetes/pkg/controllermanager/informer"
	"log"
	"sync"
)

type ControllerManager struct {
	// controller daemons
	rsController replicaset_controller.ReplicaSetController
	asController autoscaler_controller.AutoScalerController
	// informer that watch from apiserver
	podInformer informer.PodInformer
	rsInformer  informer.ReplicaSetInformer
	asInformer  informer.AutoScalerInformer
	// ensure watch order
	wg sync.WaitGroup
}

func NewControllerManager() ControllerManager {
	wg := sync.WaitGroup{}
	// informer of resources
	podInformer, _ := informer.NewPodInformer()
	rsInformer, _ := informer.NewReplicaSetInformer()
	asInformer, _ := informer.NewAutoScalerInformer()
	// controllers
	rsController, _ := replicaset_controller.NewReplicaSetController(podInformer, rsInformer, wg)
	asController, _ := autoscaler_controller.NewAutoScalerController(podInformer, rsInformer, asInformer, wg)
	return ControllerManager{
		rsController: rsController,
		asController: asController,
		podInformer:  podInformer,
		rsInformer:   rsInformer,
		asInformer:   asInformer,
		wg:           wg,
	}
}

func (cm *ControllerManager) Run() {

	// running controllers daemon
	go cm.rsController.Run()
	go cm.asController.Run()

	// informer watch must start after all controller watch
	// so we add a WaitGroup here
	cm.wg.Wait()

	// watch resources (Pod, ReplicaSet, AutoScaler) from API Server
	go cm.rsInformer.ListAndWatchReplicaSetsWithRetry()
	go cm.asInformer.ListAndWatchAutoScalersWithRetry()
	cm.podInformer.ListAndWatchPodsWithRetry()

	log.Fatalln("Unreachable here")
}
