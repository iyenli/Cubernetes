package cuberuntime

import (
	"Cubernetes/pkg/object"
	"log"
	"path/filepath"
	"strings"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func BuildPodLogsDirectory(podNamespace, podName, podUID string) string {
	return filepath.Join(podLogsRootDirectory,
		strings.Join([]string{podNamespace, podName, podUID}, "_"))
}

func toRuntimeProtocol(protocol string) runtimeapi.Protocol {
	switch protocol {
	case "TCP":
		return runtimeapi.Protocol_TCP
	case "UDP":
		return runtimeapi.Protocol_UDP
	case "SCTP":
		return runtimeapi.Protocol_SCTP
	}

	log.Printf("Unknown protocol %s, defaulting to TCP.\n", protocol)
	return runtimeapi.Protocol_TCP
}

func generateContainerMounts(pod *object.Pod, container *object.Container) []*runtimeapi.Mount {
	mounts := make([]*runtimeapi.Mount, 0)

	for _, mount := range container.VolumeMounts {
		mounts = append(mounts, &runtimeapi.Mount{
			ContainerPath: mount.MountPath,
			HostPath: findVolumeHostPath(pod, mount.Name),
		})
	}

	return mounts
}

func findVolumeHostPath(pod *object.Pod, name string) string {
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == name {
			return volume.HostPath
		}
	}
	return ""
}