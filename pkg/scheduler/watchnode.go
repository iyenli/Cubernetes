package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"time"
)

func (sr *ScheduleRuntime) WatchNode() {
	for true {
		sr.tryWatchNode()
		log.Println("[INFO]: Trying to get nodes info from apiserver after 10 secs...")
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
