package cuberuntime

import "Cubernetes/pkg/object"

const (
	CubernetesPodNameLabel       = "cubernetes.pod.name"
	CubernetesPodNameSpaceLabel  = "cubernetes.pod.namespace"
	CubernetesPodUIDLabel        = "cubernetes.pod.uid"
	CubernetesContainerNameLabel = "cubernetes.container.name"
)

func newContainerLabels(container *object.Container, pod *object.Pod) map[string]string {
	labels := map[string]string{}
	labels[CubernetesPodNameLabel] = pod.Name
	labels[CubernetesPodNameSpaceLabel] = pod.Namespace
	labels[CubernetesPodUIDLabel] = pod.UID
	labels[CubernetesContainerNameLabel] = container.Name

	return labels
}
