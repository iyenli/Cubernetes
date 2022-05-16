package autoscaler_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/informer"
	"Cubernetes/pkg/controllermanager/phase"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

const (
	minScaleInterval = time.Second * 20
	statusUpdateTime = time.Second * 8
)

type AutoScalerController interface {
	Run()
}

func NewAutoScalerController(podInformer informer.PodInformer,
	rsInformer informer.ReplicaSetInformer,
	asInformer informer.AutoScalerInformer) (AutoScalerController, error) {
	return &autoScalerController{
		podInformer: podInformer,
		rsInformer:  rsInformer,
		asInformer:  asInformer,
	}, nil
}

type autoScalerController struct {
	podInformer informer.PodInformer
	rsInformer  informer.ReplicaSetInformer
	asInformer  informer.AutoScalerInformer
	biglock     sync.Mutex
}

func (asc *autoScalerController) Run() {
	ch, cancel, err := watchobj.WatchAutoScalers()
	if err != nil {
		log.Printf("fail to watch AutoScalers from apiserver: %v\n", err)
		return
	}
	defer cancel()

	go func() {
		for {
			time.Sleep(statusUpdateTime)
			asc.updateAutoScalersRoutine()
		}
	}()

	go asc.syncLoop()

	for asEvent := range ch {
		switch asEvent.EType {
		case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
			asc.asInformer.InformAutoScaler(asEvent.AutoScaler, asEvent.EType)
		default:
			log.Fatal("[FATAL] Unknown event types: " + asEvent.EType)
		}
	}
}

func (asc *autoScalerController) syncLoop() {
	rsEventChan := asc.rsInformer.WatchRSEvent()
	defer asc.rsInformer.CloseChan(rsEventChan)

	asEventChan := asc.asInformer.WatchASEvent()
	defer asc.asInformer.CloseChan(asEventChan)

	for {
		select {
		case rsEvent := <-rsEventChan:
			asc.biglock.Lock()
			replicaSet := rsEvent.ReplicaSet
			switch rsEvent.Type {
			case types.RsCreate:
				asc.handleReplicaSetCreate(&replicaSet)
			case types.RsUpdate:
				asc.handleReplicaSetUpdate(&replicaSet)
			case types.RsRemove:
				asc.handleReplicaSetRemove(&replicaSet)
			default:
				log.Fatal("[FATAL] Unknown rsInformer event types: " + rsEvent.Type)
			}
			asc.biglock.Unlock()
		case asEvent := <-asEventChan:
			asc.biglock.Lock()
			autoScaler := asEvent.AutoScaler
			switch asEvent.Type {
			case types.AsCreate:
				asc.handleAutoScalerCreate(&autoScaler)
			case types.AsUpdate:
				asc.handleAutoScalerUpdate(&autoScaler)
			case types.AsRemove:
				asc.handleAutoScalerRemove(&autoScaler)
			default:
				log.Fatal("[FATAL] Unknown asInformer event types: " + asEvent.Type)
			}
			asc.biglock.Unlock()
		default:
			time.Sleep(time.Second * 5)
		}
	}

}

func (asc *autoScalerController) updateAutoScalersRoutine() {
	asc.biglock.Lock()
	defer asc.biglock.Unlock()

	autoScalers := asc.asInformer.ListAutoScalers()

	wg := sync.WaitGroup{}
	wg.Add(len(autoScalers))
	for _, autoScaler := range autoScalers {
		go func(as object.AutoScaler) {
			defer wg.Done()
			if as.Status != nil {
				asc.checkAndUpdateAutoScalerStatus(&as)
			}
		}(autoScaler)
	}
	wg.Wait()
}

func (asc *autoScalerController) checkAndUpdateAutoScalerStatus(as *object.AutoScaler) error {
	currentPods, err := asc.getAutoScalerPods(as)
	if err != nil {
		log.Printf("fail to get pods of AutoScaler %s: %v\n", as.Name, err)
		return err
	}

	// calculate utilazation of running pods
	runningCount := 0
	var totalCPU float64
	var totalMemory int64
	for _, pod := range currentPods {
		if phase.Running(pod.Status.Phase) {
			log.Printf("Pod %s of AutoScaler %s:\n", pod.Name, as.Name)
			log.Printf("CPU:    %f %%\n", pod.Status.ActualResourceUsage.ActualCPUUsage)
			log.Printf("Memory: %d bytes\n", pod.Status.ActualResourceUsage.ActualMemoryUsage)
			runningCount += 1
			totalCPU += pod.Status.ActualResourceUsage.ActualCPUUsage
			totalMemory += pod.Status.ActualResourceUsage.ActualMemoryUsage
		}
	}
	as.Status.ActualReplicas = runningCount
	if runningCount > 0 {
		as.Status.ActualUtilization = object.AverageUtilization{
			CPUPercentage: totalCPU / float64(runningCount),
			MemoryBytes:   totalMemory / int64(runningCount),
		}
	} else {
		log.Printf("no pod of AutoScaler %s is running, invalid utilization\n", as.Name)
		as.Status.ActualUtilization = object.AverageUtilization{}
	}

	if runningCount > 0 {
		actual := as.Status.ActualUtilization
		log.Printf("\n[AutoScaler] Average usage of %d Pods:\n", runningCount)
		log.Printf("CPU Percentage: %f %%\n", actual.CPUPercentage * 100)
		log.Printf("Memory Bytes:   %d bytes\n\n", actual.MemoryBytes)
	}

	if shouldScale, desiredReplicas := asc.computeScale(runningCount, as); shouldScale {
		rs, ok := asc.rsInformer.GetReplicaSet(as.Status.ReplicaSetUID)
		if !ok {
			log.Printf("[FATAL] lower ReplicaSet of AutoScaler %s not found\n", as.Name)
		} else {
			log.Printf("[AutoScaler] Scale desiredReplicas of %s to %d\n", as.Name, desiredReplicas)
			rs.Spec.Replicas = int32(desiredReplicas)
			as.Status.DesiredReplicas = desiredReplicas
			as.Status.LastScaleTime = time.Now()
			if _, err := crudobj.UpdateReplicaSet(*rs); err != nil {
				log.Printf("fail to update ReplicaSet Spec to apiserver\n")
			}
		}
	}

	as.Status.LastUpdateTime = time.Now()
	if _, err := crudobj.UpdateAutoScaler(*as); err != nil {
		log.Printf("fail to update autoscaler status to apiserver\n")
		return err
	}

	return nil
}

func (asc *autoScalerController) getAutoScalerPods(as *object.AutoScaler) ([]object.Pod, error) {
	return asc.podInformer.SelectPods(map[string]string{lowerReplocaSetParentUIDLabel: as.UID}), nil
}

// shouldScale, desiredReplicas
func (asc *autoScalerController) computeScale(running int, as *object.AutoScaler) (bool, int) {
	// maintain a certain time interval between 2 scales
	if time.Since(as.Status.LastScaleTime) < minScaleInterval {
		return false, -1
	}

	// not scale if running replicas doesn't match desiredReplicas
	if running != as.Status.DesiredReplicas {
		return false, -1
	}

	actual := as.Status.ActualUtilization
	limit := as.Spec.TargetUtilization

	// flag to decide the strategy of Scale:
	//  < 0: desiredReplicas -= 1
	//  = 0: no Scale
	//  > 0: desiredReplicas += 1
	scaleFlag := 0

	if limit.CPU != nil {
		if actual.CPUPercentage > limit.CPU.MaxPercentage {
			scaleFlag += 1
		} else if actual.CPUPercentage < limit.CPU.MinPercentage {
			scaleFlag -= 1
		}
	}

	if limit.Memory != nil {
		if actual.MemoryBytes > limit.Memory.MaxBytes {
			scaleFlag += 1
		} else if actual.MemoryBytes < limit.Memory.MinBytes {
			scaleFlag -= 1
		}
	}

	if scaleFlag > 0 && as.Status.DesiredReplicas < as.Spec.MaxReplicas {
		return true, as.Status.DesiredReplicas + 1
	} else if scaleFlag < 0 && as.Status.DesiredReplicas > as.Spec.MinReplicas {
		return true, as.Status.DesiredReplicas - 1
	} else {
		return false, -1
	}
}
