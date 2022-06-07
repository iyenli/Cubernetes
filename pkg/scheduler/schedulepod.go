package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"time"
)

func (sr *ScheduleRuntime) SchedulePod(pod *object.Pod) {
	if pod.Status == nil || pod.Status.NodeUID == "" {
		if pod.Status == nil {
			pod.Status = &object.PodStatus{
				ActualResourceUsage: &object.ResourceUsage{},
			}
		}

		podInfo := types.ScheduleInfo{NodeUUID: ""}
		var err error
		// Patch: easiest advanced scheduler implementation
		if len(pod.Spec.Selector) != 0 {
			nodes, err := crudobj.GetNodes()
			if err != nil {
				log.Println("[Error]: when scheduling, get nodes error:", err.Error())
				return
			}

			for _, node := range nodes {
				if object.MatchLabelSelector(pod.Spec.Selector, node.Labels) {
					podInfo = types.ScheduleInfo{NodeUUID: node.UID}
					break
				}
			}
		} else {
			podInfo, err = sr.Implement.Schedule()
			if err != nil {
				log.Println("[Error]: when scheduling, error:", err.Error())
				return
			}
		}

		if podInfo.NodeUUID == "" {
			log.Println("[Warn]: No node to schedule this node")
		}

		err = sr.SendPodScheduleInfoBack(pod, &podInfo)
		if err != nil {
			log.Println("[Error]: when sending scheduler result,", err.Error())
		}
	}
}

func (sr *ScheduleRuntime) SendPodScheduleInfoBack(podToSchedule *object.Pod, info *types.ScheduleInfo) error {
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

func (sr *ScheduleRuntime) WatchPod() {
	for {
		sr.tryWatchPod()
		time.Sleep(WatchRetryIntervalSec * time.Second)
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
		log.Printf("[Error]: Error occurs when watching pods: %v", err)
		return
	}
	defer cancel()

	for {
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
					log.Println("[Info]: Delete Pod, do nothing")
				default:
					log.Panic("[Fatal]: Unsupported types in watching pod.")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
