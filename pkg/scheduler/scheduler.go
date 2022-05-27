package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/scheduler/RR"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"sync"
)

const WatchRetryIntervalSec = 10

type ScheduleRuntime struct {
	Implement types.Scheduler
}

func NewScheduler() *ScheduleRuntime {
	scheduler := RR.SchedulerRR{
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
	// Init sr first
	sr.Init()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		sr.WatchNode()
	}()

	go func() {
		defer wg.Done()
		sr.WatchPod()
	}()

	go func() {
		defer wg.Done()
		sr.WatchJob()
	}()

	wg.Wait()
	log.Fatalln("Unreachable here")
}

func (sr *ScheduleRuntime) Init() {
	log.Println("[INFO]: Init scheduler using exist jobs and pods...")
	pods, err := crudobj.GetPods()
	if err != nil {
		log.Println("[Error]: when init scheduler and ask for pods,", err.Error())
	}
	for _, pod := range pods {
		sr.SchedulePod(&pod)
	}

	jobs, err := crudobj.GetGpuJobs()
	if err != nil {
		log.Println("[Error]: when init scheduler and ask for jobs", err.Error())
	}
	for _, job := range jobs {
		sr.ScheduleJob(&job)
	}
}
