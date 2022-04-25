package cuberuntime

import (
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
