package Advanced

import (
	"Cubernetes/pkg/scheduler/types"
)

//var ErrNoNodesToSchedule = errors.New("not any nodes for scheduling")

type SchedulerAdvanced struct {
	NumOfNodes  int32
	NameOfNodes []string
	// TODO: Implement it!
}

func (rr *SchedulerAdvanced) Init() error {
	rr.NumOfNodes = 0
	rr.NameOfNodes = make([]string, 0)

	return nil
}

func (rr *SchedulerAdvanced) AddNode(info *types.NodeInfo) error {
	// No redundant node:)
	return nil
}

func (rr *SchedulerAdvanced) RemoveNode(info *types.NodeInfo) error {
	if rr.NumOfNodes == 0 {
		return nil
	}

	return nil
}

func (rr *SchedulerAdvanced) Schedule() (types.PodInfo, error) {
	return types.PodInfo{NodeUUID: ""}, nil
}
