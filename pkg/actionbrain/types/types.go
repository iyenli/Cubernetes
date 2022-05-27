package types

import "Cubernetes/pkg/object"

type ActionEventType string

const (
	ActionCreate ActionEventType = "create"
	ActionRemove ActionEventType = "remove"
)

type ActionEvent struct {
	Type   ActionEventType
	Action object.Action
}

type ActorEventType string

const (
	ActorCreate ActorEventType = "create"
	ActorRemove ActorEventType = "remove"
)

type ActorEvent struct {
	Type  ActorEventType
	Actor object.Actor
}
