package types

import "Cubernetes/pkg/object"

type IngressEventType string

const (
	IngressCreate IngressEventType = "CreateIngress"
	IngressUpdate IngressEventType = "UpdateIngress"
	IngressRemove IngressEventType = "RemoveIngress"
)

type IngressEvent struct {
	Type    IngressEventType
	Ingress object.Ingress
}
