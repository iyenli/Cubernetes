package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"time"
)

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

func (sr *ScheduleRuntime) WatchJob() {
	for true {
		sr.tryWatchJob()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}
