package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
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

func (sr *ScheduleRuntime) SchedulePod(pod *object.Pod) {
	if pod.Status == nil || pod.Status.NodeUID == "" {
		if pod.Status == nil {
			pod.Status = &object.PodStatus{
				ActualResourceUsage: &object.ResourceUsage{},
			}
		}

		podInfo, err := sr.Implement.Schedule()
		if err != nil {
			log.Println("[Error]: when scheduling, error:", err.Error())
		}

		err = sr.SendPodScheduleInfoBack(pod, &podInfo)
		if err != nil {
			log.Println("[Error]: when sending scheduler result,", err.Error())
		}
	}
}

func (sr *ScheduleRuntime) ScheduleJob(job *object.GpuJob) {
	// only support job has checked files and never be scheduled
	if job.Status.NodeUID == "" && job.Status.Phase == object.JobCreated {
		podInfo, err := sr.Implement.Schedule()
		if err != nil {
			log.Println("[Error]: when scheduling, error:", err.Error())
		}
		err = sr.SendJobScheduleInfoBack(job, &podInfo)
		if err != nil {
			log.Println("[Error]: when sending scheduler result,", err.Error())
		}
	}
}

func (sr *ScheduleRuntime) SendPodScheduleInfoBack(podToSchedule *object.Pod, info *types.PodInfo) error {
	podToSchedule.Status.NodeUID = info.NodeUUID
	podToSchedule.Status.Phase = object.PodBound

	_, err := crudobj.UpdatePod(*podToSchedule)
	if err != nil {
		log.Println("[INFO]: Update pod failed")
		return err
	}

	log.Println("[INFO]: Schedule pod", podToSchedule.UID, "to node", info.NodeUUID)
	return nil
}

func (sr *ScheduleRuntime) SendJobScheduleInfoBack(jobToSchedule *object.GpuJob, info *types.PodInfo) error {
	jobToSchedule.Status.NodeUID = info.NodeUUID

	_, err := crudobj.UpdateGpuJob(*jobToSchedule)
	if err != nil {
		log.Println("[INFO]: Update pod failed")
		return err
	}

	log.Println("[INFO]: Schedule job", jobToSchedule.UID, "to node", info.NodeUUID)
	return nil
}

func (sr *ScheduleRuntime) WatchNode() {
	for true {
		sr.tryWatchNode()
		log.Println("[INFO]: Trying to get nodes info from apiserver after 10 secs...")
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (sr *ScheduleRuntime) WatchPod() {
	for true {
		sr.tryWatchPod()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (sr *ScheduleRuntime) WatchJob() {
	for true {
		sr.tryWatchJob()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (sr *ScheduleRuntime) tryWatchNode() {
	if allNodes, err := crudobj.GetNodes(); err != nil {
		log.Printf("[INFO]: fail to get all nodes from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, node := range allNodes {
			_ = sr.Implement.AddNode(&types.NodeInfo{NodeUUID: node.UID})
		}
	}

	ch, handler, err := watchobj.WatchNodes()
	if err != nil {
		log.Println("[INFO]: Get nodes channel failed")
		return
	}
	defer handler()

	for true {
		select {
		case nodeEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			} else {
				if nodeEvent.EType == watchobj.EVENT_PUT {
					if nodeEvent.Node.Status == nil {
						continue
					}
					if nodeEvent.Node.Status.Condition.Ready == false {
						log.Println("[INFO]: Scheduler may removed a node: ", nodeEvent.Node.UID)
						err := sr.Implement.RemoveNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
						if err != nil {
							log.Println("[Error]: remove node failed")
						}
					}
					if nodeEvent.Node.Status.Condition.Ready == true {
						log.Println("[INFO]: Scheduler may added a node: ", nodeEvent.Node.UID)
						err := sr.Implement.AddNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
						if err != nil {
							log.Println("[error]: add node failed")
						}
					}
				} else if nodeEvent.EType == watchobj.EVENT_DELETE {
					log.Println("[INFO]: Scheduler may removed a node: ", nodeEvent.Node.UID)
					err := sr.Implement.RemoveNode(&types.NodeInfo{NodeUUID: nodeEvent.Node.UID})
					if err != nil {
						log.Println("[error]: remove node failed")
					}
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (sr *ScheduleRuntime) tryWatchPod() {
	if allPods, err := crudobj.GetPods(); err != nil {
		log.Printf("[INFO]: fail to get all pods from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, pod := range allPods {
			sr.SchedulePod(&pod)
		}
	}

	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Printf("Error occurs when watching pods: %v", err)
		return
	}
	defer cancel()

	for true {
		select {
		case podEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			} else {
				switch podEvent.EType {
				case watchobj.EVENT_PUT:
					sr.SchedulePod(&podEvent.Pod)
				case watchobj.EVENT_DELETE:
					log.Println("[Info]: Delete pod, do nothing")
				default:
					log.Panic("[Fatal]: Unsupported types in watching pod.")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (sr *ScheduleRuntime) tryWatchJob() {
	if allJobs, err := crudobj.GetGpuJobs(); err != nil {
		log.Printf("[INFO]: fail to get all jobs from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, job := range allJobs {
			sr.ScheduleJob(&job)
		}
	}

	ch, cancel, err := watchobj.WatchGpuJobs()
	if err != nil {
		log.Printf("Error occurs when watching jobs: %v", err)
		return
	}
	defer cancel()

	for true {
		select {
		case jobEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			} else {
				switch jobEvent.EType {
				case watchobj.EVENT_PUT:
					sr.ScheduleJob(&jobEvent.GpuJob)
				case watchobj.EVENT_DELETE:
					log.Println("[Info]: delete pod, do nothing")
				default:
					log.Panic("[Fatal]: Unsupported types in watching pod")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
