package cuberuntime

import (
	"Cubernetes/pkg/object"
	"log"
	"path/filepath"
	"strings"
)

func BuildPodLogsDirectory(podNamespace, podName, podUID string) string {
	return filepath.Join(podLogsRootDirectory,
		strings.Join([]string{podNamespace, podName, podUID}, "_"))
}

func toPortProtocol(protocol string) string {
	switch protocol {
	case "TCP":
		return "/tcp"
	case "UDP":
		return "/udp"
	case "SCTP":
		return "/sctp"
	}

	log.Printf("Unknown protocol%s, defaulting to TCP.\n", protocol)
	return "/tcp"
}

func findVolumeHostPath(pod *object.Pod, name string) string {
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == name {
			return volume.HostPath
		}
	}
	return ""
}

func buildLabelSelector(label, value string) string {
	return strings.Join([]string{label, value}, "=")
}
