package cuberuntime

import "Cubernetes/pkg/object"

const (
	PodNameLabel       = "cubernetes.pod.name"
	PodNameSpaceLabel  = "cubernetes.pod.namespace"
	PodUIDLabel        = "cubernetes.pod.uid"
	ContainerNameLabel = "cubernetes.container.name"
	ContainerTypeLabel = "cubernetes.container.types"

	ContainerTypeContainer = "container"
	ContainerTypeSandbox   = "sandbox"
)

func newContainerLabels(container *object.Container, pod *object.Pod) map[string]string {
	labels := map[string]string{}
	labels[PodNameLabel] = pod.Name
	labels[PodNameSpaceLabel] = pod.Namespace
	labels[PodUIDLabel] = pod.UID

	labels[ContainerNameLabel] = container.Name
	labels[ContainerTypeLabel] = ContainerTypeContainer

	return labels
}

func newSandboxLabels(pod *object.Pod) map[string]string {
	labels := map[string]string{}

	for k, v := range pod.Labels {
		labels[k] = v
	}
	labels[PodNameLabel] = pod.Name
	labels[PodNameSpaceLabel] = pod.Namespace
	labels[PodUIDLabel] = pod.UID

	labels[ContainerTypeLabel] = ContainerTypeSandbox

	return labels
}
