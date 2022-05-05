package watchobj

import "Cubernetes/pkg/object"

type EventType string

const (
	EVENT_PUT    EventType = "PUT"
	EVENT_DELETE EventType = "DELETE"
)

const MSG_DELIM byte = 26
const WATCH_CONFIRM string = "watch started"

type ObjEvent struct {
	EType  EventType `json:"eType"`
	Path   string    `json:"path"`
	Object string    `json:"object"`
}

type PodEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Pod will only have its UID
	Pod object.Pod
}

type ServiceEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// Service will only have its UID
	Service object.Service
}

type ReplicaSetEvent struct {
	EType EventType
	// if EType == EVENT_DELETE,
	// ReplicaSet will only have its UID
	ReplicaSet object.ReplicaSet
}
