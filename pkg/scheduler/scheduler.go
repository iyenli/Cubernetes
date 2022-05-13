package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/scheduler/RR"
	"Cubernetes/pkg/scheduler/types"
	"log"
)

type ScheduleRuntime struct {
	Implement types.Scheduler
}

// RealScheduler Choose one: RR.SchedulerRR / Advanced.SchedulerAdvanced
type RealScheduler = RR.SchedulerRR

func NewScheduler() *ScheduleRuntime {
	scheduler := RealScheduler{
		NumOfNodes:  0,
		NameOfNodes: []string{},
	}
	err := scheduler.Init()
	if err != nil {
		log.Panicln("[Panic]: Error when init scheduler")
		return nil
	}

	return &ScheduleRuntime{
		Implement: &scheduler,
	}
}

func (sr *ScheduleRuntime) Run() {
	ch, cancel, err := watchobj.WatchPods()

	if err != nil {
		log.Panicf("Error occurs when watching pods: %v", err)
	}
	defer cancel()

	go func() {
		err := sr.WatchNode()
		if err != nil {
			return
		}
	}()

	for podEvent := range ch {
		switch podEvent.EType {
		case watchobj.EVENT_PUT:
			pod := podEvent.Pod
			if pod.Status == nil || pod.Status.PodUID == "" {
				if pod.Status == nil {
					pod.Status = &object.PodStatus{
						ActualResourceUsage: &object.ResourceUsage{},
					}
				}

				podInfo, err := sr.Implement.Schedule()
				if err != nil {
					log.Println("Error happened when scheduling")
				}

				err = sr.SendScheduleInfoBack(&pod, &podInfo)
				if err != nil {
					log.Println("Error happened when sending scheduler result")
				}
			}

		case watchobj.EVENT_DELETE:
			log.Println("[Info]: Delete pod, do nothing")
		default:
			log.Panic("Unsupported types in watch pod.")
		}
	}

	log.Fatalln("Unreachable here")
}

func (sr *ScheduleRuntime) SendScheduleInfoBack(podToSchedule *object.Pod, info *types.PodInfo) error {
	podToSchedule.Status.PodUID = info.NodeUUID

	_, err := crudobj.UpdatePod(*podToSchedule)
	log.Println("[INFO]: Schedule pod ", podToSchedule.UID, "To node ", info.NodeUUID)

	if err != nil {
		log.Println("Update pod failed")
		return err
	}

	return nil
}

func (sr *ScheduleRuntime) WatchNode() error {
	// init: Get existed scheduler
	ch, handler, err := watchobj.WatchNodes()
	if err != nil {
		log.Fatalln("[Fatal]: Get nodes to init failed")
	}

	defer handler()
	for nodeEvent := range ch {
		if nodeEvent.EType == watchobj.EVENT_PUT {
			if nodeEvent.Node.Status == nil {
				continue
			}

			if nodeEvent.Node.Status.Condition.Ready == false {
				log.Println("[INFO] Scheduler may removed a node: ", nodeEvent.Node.UID)
				err := sr.Implement.RemoveNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
				if err != nil {
					log.Println("[error]: remove node failed")
				}
			}
			if nodeEvent.Node.Status.Condition.Ready == true {
				log.Println("[INFO] Scheduler may added a node: ", nodeEvent.Node.UID)
				err := sr.Implement.AddNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
				if err != nil {
					log.Println("[error]: add node failed")
				}
			}
		} else if nodeEvent.EType == watchobj.EVENT_DELETE {
			log.Println("[INFO] Scheduler may removed a node: ", nodeEvent.Node.UID)
			err := sr.Implement.RemoveNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
			if err != nil {
				log.Println("[error]: remove node failed")
			}
		}
	}

	return nil
}
