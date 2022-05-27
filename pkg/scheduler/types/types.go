package types

type NodeInfo struct {
	NodeUUID string
}

type ScheduleInfo struct {
	NodeUUID string
}

type Scheduler interface {
	Init() error

	AddNode(Info *NodeInfo) error

	RemoveNode(Info *NodeInfo) error

	Schedule() (ScheduleInfo, error)
}
