package cuberuntime

import (
	"Cubernetes/pkg/object"
	"fmt"
	"log"
	"net"
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

// determinePodSandboxIP determines the IP addresses of the given pod sandbox.
func (m *cubeRuntimeManager) determinePodSandboxIPs(podNamespace, podName string, podSandbox *runtimeapi.PodSandboxStatus) []string {
	podIPs := make([]string, 0)
	if podSandbox.Network == nil {
		log.Printf("Pod %s's Sandbox status doesn't have network information, cannot report IPs", podName)
		return podIPs
	}

	// ip could be an empty string if runtime is not responsible for the
	// IP (e.g., host networking).

	// pick primary IP
	if len(podSandbox.Network.Ip) != 0 {
		if net.ParseIP(podSandbox.Network.Ip) == nil {
			log.Printf("Pod %s's Sandbox reported an unparseable primary IP %s", podName, podSandbox.Network.Ip)
			return nil
		}
		podIPs = append(podIPs, podSandbox.Network.Ip)
	}

	// pick additional ips, if cri reported them
	for _, podIP := range podSandbox.Network.AdditionalIps {
		if nil == net.ParseIP(podIP.Ip) {
			log.Printf("Pod %s's Sandbox reported an unparseable additional IP %s", podName, podIP.Ip)
			return nil
		}
		podIPs = append(podIPs, podIP.Ip)
	}

	return podIPs
}
