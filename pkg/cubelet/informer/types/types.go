package types

import "Cubernetes/pkg/object"

type EventType string

const (
	Create EventType = "create"
	Update EventType = "update"
	Remove EventType = "remove"
)

type PodEvent struct {
	Type EventType
	Pod  object.Pod
}

type JobEvent struct {
	Type EventType
	Job  object.GpuJob
}
