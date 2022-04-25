package cuberuntime

import (
	"Cubernetes/pkg/object"
	"fmt"
	"log"
	"os"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (m *cubeRuntimeManager) createPodSandbox(pod *object.Pod, attempt uint32) (string, string, error) {
	podSandboxConfig, err := m.generatePodSandboxConfig(pod, attempt)
	if err != nil {
		message := fmt.Sprintf("Failed to generate sandbox config for pod %s: %v", pod.Name, err)
		log.Println(message)
		return "", message, err
	}

	err = os.MkdirAll(podSandboxConfig.LogDirectory, 0755)
	if err != nil {
		message := fmt.Sprintf("Failed to create log directory for pod %s: %v", pod.Name, err)
		log.Println(message)
		return "", message, err
	}

	// use default runtime handler now
	runtimeHandler := ""

	podSandBoxID, err := m.runtimeService.RunPodSandbox(podSandboxConfig, runtimeHandler)
	if err != nil {
		message := fmt.Sprintf("Failed to create sandbox for pod %s: %v", pod.Name, err)
		log.Println(message)
		return "", message, err
	}

	return podSandBoxID, "", nil
}

func (m *cubeRuntimeManager) generatePodSandboxConfig(pod *object.Pod, attempt uint32) (*runtimeapi.PodSandboxConfig, error) {
	podSandboxConfig := &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Uid:       pod.UID,
			Attempt:   attempt,
		},
		LogDirectory: BuildPodLogsDirectory(pod.Namespace, pod.Name, pod.UID),
		Labels:       newPodLabels(pod),
		Annotations:  pod.Annotations,
	}

	// configure port mapping
	portMappings := []*runtimeapi.PortMapping{}
	for _, c := range pod.Spec.Containers {
		for _, p := range c.Ports {
			portMappings = append(portMappings, &runtimeapi.PortMapping{
				HostIp:        p.HostIP,
				HostPort:      p.HostPort,
				ContainerPort: p.ContainerPort,
				Protocol:      toRuntimeProtocol(p.Protocol),
			})
		}
	}
	if len(portMappings) > 0 {
		podSandboxConfig.PortMappings = portMappings
	}

	// TODO: configure pod request & limit resources

	return podSandboxConfig, nil
}
