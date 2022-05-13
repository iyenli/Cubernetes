package types

import "Cubernetes/pkg/object"

type PodEventType string
type ServiceEventType string

const (
	PodCreate PodEventType = "CreatePod"
	PodUpdate PodEventType = "UpdatePod"
	PodRemove PodEventType = "RemovePod"
)

const (
	ServiceCreate ServiceEventType = "CreateService"
	ServiceUpdate ServiceEventType = "UpdateService"
	ServiceRemove ServiceEventType = "RemoveService"
)

type PodEvent struct {
	Type PodEventType
	Pod  object.Pod
}

type ServiceEvent struct {
	Type    ServiceEventType
	Service object.Service
}
