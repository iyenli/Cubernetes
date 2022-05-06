package types

import "Cubernetes/pkg/object"

type PodEventType string

const (
	PodCreate PodEventType = "create"
	PodUpdate PodEventType = "update"
	PodRemove PodEventType = "remove"
)

type PodEvent struct {
	Type PodEventType
	Pod  *object.Pod
}
