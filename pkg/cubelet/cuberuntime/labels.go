package cuberuntime

import "Cubernetes/pkg/object"

const (
	CubernetesPodNameLabel       = "cubernetes.pod.name"
	CubernetesPodNameSpaceLabel  = "cubernetes.pod.namespace"
	CubernetesPodUIDLabel        = "cubernetes.pod.uid"
	CubernetesContainerNameLabel = "cubernetes.container.name"
	CubernetesContainerTypeLabel = "cubernetes.container.type"

	ContainerTypeContainer = "container"
	ContainerTypeSandbox   = "sandbox"
)

func newContainerLabels(container *object.Container, pod *object.Pod) map[string]string {
	labels := map[string]string{}
	labels[CubernetesPodNameLabel] = pod.Name
	labels[CubernetesPodNameSpaceLabel] = pod.Namespace
	labels[CubernetesPodUIDLabel] = pod.UID

	labels[CubernetesContainerNameLabel] = container.Name
	labels[CubernetesContainerTypeLabel] = ContainerTypeContainer

	return labels
}

func newSandboxLabels(pod *object.Pod) map[string]string {
	labels := map[string]string{}

	for k, v := range pod.Labels {
		labels[k] = v
	}
	labels[CubernetesPodNameLabel] = pod.Name
	labels[CubernetesPodNameSpaceLabel] = pod.Namespace
	labels[CubernetesPodUIDLabel] = pod.UID

	labels[CubernetesContainerTypeLabel] = ContainerTypeSandbox

	return labels
}
