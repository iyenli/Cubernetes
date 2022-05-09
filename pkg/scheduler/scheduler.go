package scheduler

import (
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
		log.Panicln("Error when init scheduler")
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
			podInfo, err := sr.Implement.Schedule()
			if err != nil {
				log.Println("Error happened when scheduling")
			}

			err = sr.SendScheduleInfoBack(&podEvent.Pod, &podInfo)
			if err != nil {
				log.Println("Error happened when sending schedluer result")
			}
		case watchobj.EVENT_DELETE:
			log.Println("Delete pod, do nothing")
		default:
			log.Panic("Unsupported type in watch pod.")
		}
	}

	log.Fatalln("Unreachable here")
}

func (sr *ScheduleRuntime) SendScheduleInfoBack(podToSchedule *object.Pod, info *types.PodInfo) error {
	// TODO: Write Schedule info back
	return nil
}

func (sr *ScheduleRuntime) WatchNode() error {
	// TODO
	for {
		return nil
	}
	log.Panicf("Unreachable here")
	return nil
}
