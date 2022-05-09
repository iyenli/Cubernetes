package RR

import (
	"Cubernetes/pkg/scheduler/types"
	"errors"
	"sync/atomic"
)

var ErrNoNodesToSchedule = errors.New("not any nodes for scheduling")

type SchedulerRR struct {
	NumOfNodes  int32
	NameOfNodes []string

	Next int32
}

func (rr *SchedulerRR) Init() error {
	rr.Next = 0
	rr.NumOfNodes = 0
	rr.NameOfNodes = make([]string, 0)

	return nil
}

func (rr *SchedulerRR) AddNode(info *types.NodeInfo) error {
	// No redundant node:)
	for _, node := range rr.NameOfNodes {
		if node == info.NodeUUID {
			return nil
		}
	}

	atomic.AddInt32(&rr.NumOfNodes, 1)
	rr.NameOfNodes = append(rr.NameOfNodes, info.NodeUUID)

	return nil
}

func (rr *SchedulerRR) RemoveNode(info *types.NodeInfo) error {
	if rr.NumOfNodes == 0 {
		return nil
	}

	for idx, node := range rr.NameOfNodes {
		if node == info.NodeUUID {
			rr.NameOfNodes = append(rr.NameOfNodes[:idx], rr.NameOfNodes[idx+1:]...)
			atomic.AddInt32(&rr.NumOfNodes, -1)
			return nil
		}
	}

	return nil
}

func (rr *SchedulerRR) Schedule() (types.PodInfo, error) {
	if rr.NumOfNodes == 0 {
		return types.PodInfo{NodeUUID: ""}, ErrNoNodesToSchedule
	}

	n := atomic.AddInt32(&rr.Next, 1)
	return types.PodInfo{NodeUUID: rr.NameOfNodes[((n - 1) % rr.NumOfNodes)]}, nil
}
