package dockershim

import (
	"Cubernetes/pkg/object"
	"strings"
)

const (
	cubePrefix           = "c8s"
	sandboxContainerName = "PODSandbox"
	nameDelimiter        = "_"
)

func MakeSandboxName(pod *object.Pod) string {
	return strings.Join([]string{
		cubePrefix,
		pod.Name,
		sandboxContainerName,
		pod.UID,
	}, nameDelimiter)
}

func MakeContainerName(pod *object.Pod, container *object.Container) string {
	return strings.Join([]string{
		cubePrefix,
		pod.Name,
		container.Name,
		pod.UID,
	}, nameDelimiter)
}

func ParseSandboxName(sandboxName string) string {
	name := strings.Trim(sandboxName, "/")
	return strings.Split(name, nameDelimiter)[1]
}

func ParseContainerName(containerName string) string {
	name := strings.Trim(containerName, "/")
	return strings.Split(name, nameDelimiter)[2]
}
