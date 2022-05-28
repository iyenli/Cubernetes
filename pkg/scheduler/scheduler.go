package scheduler

import (
	"Cubernetes/pkg/scheduler/RR"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"sync"
	"time"
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
	log.Println("[INFO]: Init Scheduler with current nodes, it may take 2 seconds...")

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		sr.WatchNode()
	}()

	// Get all nodes in sr
	time.Sleep(2 * time.Second)

	go func() {
		defer wg.Done()
		sr.WatchPod()
	}()

	go func() {
		defer wg.Done()
		sr.WatchJob()
	}()

	go func() {
		defer wg.Done()
		sr.WatchActor()
	}()

	wg.Wait()
	log.Fatalln("Unreachable here")
}
