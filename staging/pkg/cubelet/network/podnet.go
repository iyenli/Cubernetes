package network

import (
	"Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/object"
	"fmt"
	gocni "github.com/containerd/go-cni"
	"log"
	"net"
	osexec "os/exec"
	"strings"
)

const (
	defaultNetworkNamespace = "/var/run/netns/default"
)

func InitPodNetwork(cni gocni.CNI, podStatus *container.PodStatus) error {
	// TODO: Run pause docker and add it to SandboxStatuses

	if len(podStatus.SandboxStatuses) < 1 {
		log.Panicln("Error: Pause sandbox create failed")
	}
	if podStatus.NetworkNamespace == "" {
		podStatus.NetworkNamespace = defaultNetworkNamespace
	}

	result, err := SetUpPod(cni, podStatus.NetworkNamespace, podStatus.UID, container.ContainerID{ID: podStatus.SandboxStatuses[0].Id})
	if err != nil {
		log.Println("Setup pod failed.")
		return err
	}

	podStatus.PodNetWork.IP = net.ParseIP(result.Interfaces["eth"].IPConfigs[0].IP.String())
	return nil
}

func ReleaseNetwork(cni gocni.CNI, podStatus *container.PodStatus) error {
	err := TearDownPod(cni, podStatus.NetworkNamespace, podStatus.UID, container.ContainerID{ID: podStatus.SandboxStatuses[0].Id})
	if err != nil {
		log.Println("Teardown pod failed.")
		return err
	}
	return nil
}

// ConstructPodPortMapping creates a PodPortMapping from the ports specified in the pod's
// containers.
func ConstructPodPortMapping(pod *object.Pod, podIP net.IP) *PodPortMapping {
	portMappings := make([]*PortMapping, 0)
	for _, c := range pod.Spec.Containers {
		for _, port := range c.Ports {
			portMappings = append(portMappings, &PortMapping{
				Name:          port.Name,
				HostPort:      port.HostPort,
				ContainerPort: port.ContainerPort,
				Protocol:      port.Protocol,
				HostIP:        port.HostIP,
			})
		}
	}

	return &PodPortMapping{
		Namespace:    pod.Namespace,
		Name:         pod.Name,
		PortMappings: portMappings,
		IP:           podIP,
	}
}

func GetPodIP(nsenterPath, netnsPath, interfaceName string) (net.IP, error) {
	// Only support IPv4 for simplicity
	ip, err := getOnePodIP(nsenterPath, netnsPath, interfaceName, "-4")
	if err != nil {
		return nil, err
	}

	return ip, nil
}

func getOnePodIP(nsenterPath, netnsPath, interfaceName, addrType string) (net.IP, error) {
	// Try to retrieve ip inside container network namespace
	cmd := osexec.Command(nsenterPath, fmt.Sprintf("--net=%s", netnsPath), "-F", "--",
		"ip", "-o", addrType, "addr", "show", "dev", interfaceName, "scope", "global")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("unexpected command output %s with error: %v", output, err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected command output %s", output)
	}

	fields := strings.Fields(lines[0])
	if len(fields) < 4 {
		return nil, fmt.Errorf("unexpected address output %s ", lines[0])
	}

	ip, _, err := net.ParseCIDR(fields[3])
	if err != nil {
		return nil, fmt.Errorf("cni failed to parse ip from output %s due to %v", output, err)
	}

	return ip, nil
}
