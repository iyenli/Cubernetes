package types

import "Cubernetes/pkg/object"

type ActionEventType string

const (
	ActionCreate ActionEventType = "create"
	ActionUpdate ActionEventType = "update"
	ActionRemove ActionEventType = "remove"
)

type ActionEvent struct {
	Type   ActionEventType
	Action object.Action
}

type ActorEventType string

const (
	ActorCreate ActorEventType = "create"
	ActorUpdate ActorEventType = "update"
	ActorRemove ActorEventType = "remove"
)

type ActorEvent struct {
	Type  ActorEventType
	Actor object.Actor
}
