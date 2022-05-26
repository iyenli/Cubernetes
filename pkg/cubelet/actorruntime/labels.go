package actor_runtime

import (
	"Cubernetes/pkg/object"
	"strings"
)

const (
	ActorNameLabel     = "cubernetes.actor.name"
	ActorUIDLabel      = "cubernetes.actor.uid"
	ContainerTypeLabel = "cubernetes.actor.container.types"

	ContainerTypeContainer = "container"
	ContainerTypeSandbox   = "sandbox"
)

func newContainerLabels(actor *object.Actor) map[string]string {
	labels := map[string]string{
		ActorUIDLabel:      actor.UID,
		ActorNameLabel:     actor.Name,
		ContainerTypeLabel: ContainerTypeContainer,
	}

	for k, v := range actor.Labels {
		labels[k] = v
	}

	return labels
}

func newSandboxLabels(actor *object.Actor) map[string]string {
	labels := map[string]string{
		ActorUIDLabel:      actor.UID,
		ActorNameLabel:     actor.Name,
		ContainerTypeLabel: ContainerTypeSandbox,
	}

	for k, v := range actor.Labels {
		labels[k] = v
	}

	return labels
}

func buildLabelSelector(label, value string) string {
	return strings.Join([]string{label, value}, "=")
}
