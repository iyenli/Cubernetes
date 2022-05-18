package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type JobInformer interface {
	ListAndWatchJobsWithRetry()
	WatchJobEvent() <-chan types.JobEvent
	SetNodeUID(uid string)
}

type cubeJobInformer struct {
	jobEvent chan types.JobEvent
	jobCache map[string]bool

	nodeUID string
}

func NewJobInformer() (JobInformer, error) {
	return &cubeJobInformer{
		jobEvent: make(chan types.JobEvent),
		jobCache: make(map[string]bool),
	}, nil
}

func (c *cubeJobInformer) SetNodeUID(uid string) {
	if c.nodeUID != "" {
		log.Printf("[FATAL]: Node ID already set!\n")
	}
	c.nodeUID = uid
}

func (c *cubeJobInformer) ListAndWatchJobsWithRetry() {
	defer close(c.jobEvent)
	for true {
		c.tryListAndWatchJobs()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (c *cubeJobInformer) tryListAndWatchJobs() {
	if allJobs, err := crudobj.GetGpuJobs(); err != nil {
		log.Printf("[INFO]: fail to get all jobs from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, job := range allJobs {
			if job.Status.NodeUID == c.nodeUID && job.Status.Phase == object.JobCreated {
				err := c.informJob(job, watchobj.EVENT_PUT)
				if err != nil {
					log.Printf("[INFO]: Inform job failed when initializing job informer")
				} else {
					c.jobCache[job.UID] = true
				}
			}
		}
	}

	ch, cancel, err := watchobj.WatchGpuJobs()
	if err != nil {
		log.Printf("fail to watch pods from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for true {
		select {
		case jobEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			} else if jobEvent.GpuJob.Status.Phase != object.JobCreated {
				log.Printf("[INFO]: Job received, phase is %v. Job uuid is %v \n", jobEvent.GpuJob.Status.Phase, jobEvent.GpuJob.UID)
				log.Printf("[INFO]: Not handle it\n")
				continue
			} else if jobEvent.GpuJob.Status.NodeUID != c.nodeUID && jobEvent.EType != watchobj.EVENT_DELETE {
				log.Printf("[INFO]: Job received, Not my job. Job uuid is %v \n", jobEvent.GpuJob.UID)
				continue
			} else if jobEvent.EType == watchobj.EVENT_PUT || jobEvent.EType == watchobj.EVENT_DELETE {
				log.Printf("[INFO]: My job received, Job uuid is %v \n", jobEvent.GpuJob.UID)
				err := c.informJob(jobEvent.GpuJob, jobEvent.EType)
				if err != nil {
					return
				}
			} else {
				log.Panic("[Error]: Unsupported types in watch pod.")
			}
		default:
			time.Sleep(time.Second)
		}
	}

}

func (c *cubeJobInformer) WatchJobEvent() <-chan types.JobEvent {
	return c.jobEvent
}

func (c *cubeJobInformer) informJob(newJob object.GpuJob, eType watchobj.EventType) error {
	_, exist := c.jobCache[newJob.UID]
	if eType == watchobj.EVENT_DELETE {
		log.Printf("[INFO]: Delete job request received, I don't know what todo")
		return nil
	} else {
		if exist {
			log.Printf("[INFO]: pod %s has been handled before\n", newJob.UID)
		} else {
			c.jobEvent <- types.JobEvent{
				Type: types.Create,
				Job:  newJob,
			}
		}
	}
	return nil
}
