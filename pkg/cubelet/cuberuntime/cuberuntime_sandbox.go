package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/object"
	"log"
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockernat "github.com/docker/go-connections/nat"
)

const (
	defaultPauseImage = "docker/desktop-kubernetes-pause:3.7"
)

// sandboxName, sandboxID, err
func (m *cubeRuntimeManager) createPodSandbox(pod *object.Pod) (string, string, error) {
	err := m.dockerRuntime.PullImage(defaultPauseImage)
	if err != nil {
		log.Printf("ensure image for sandbox of pod %s failed\n", pod.Name)
		return "", "", err
	}

	podSandboxConfig := generatePodSandboxConfig(pod)
	sandboxID, err := m.dockerRuntime.CreateContainer(podSandboxConfig)
	if err != nil {
		log.Printf("fail to create sandbox of pod %s\n", pod.Name)
		return "", "", err
	}

	err = m.dockerRuntime.StartContainer(sandboxID)
	if err != nil {
		log.Printf("fail to start sandbox of pod %s\n", pod.Name)
		return "", "", err
	}

	return podSandboxConfig.Name, sandboxID, nil
}

func (m *cubeRuntimeManager) getSandboxStatusesByPodUID(UID string) ([]*cubecontainer.SandboxStatus, error) {
	filter := dockertypes.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", buildLabelSelector(ContainerTypeLabel, ContainerTypeSandbox)),
			filters.Arg("label", buildLabelSelector(PodUIDLabel, UID)),
		),
	}

	sandboxes, err := m.dockerRuntime.ListContainers(filter)
	if err != nil {
		log.Printf("fail to list pod sandbox %s: %v\n", UID, err)
		return nil, err
	}

	if len(sandboxes) == 0 {
		return nil, nil
	}

	statuses := make([]*cubecontainer.SandboxStatus, len(sandboxes))
	for i, sandbox := range sandboxes {
		statuses[i] = toSandboxStatus(&sandbox)
	}

	return statuses, nil
}

func generatePodSandboxConfig(pod *object.Pod) *dockertypes.ContainerCreateConfig {
	sandboxName := dockershim.MakeSandboxName(pod)

	exposedPorts := dockernat.PortSet{}
	portBindings := map[dockernat.Port][]dockernat.PortBinding{}

	// port bindings
	for _, c := range pod.Spec.Containers {
		for _, p := range c.Ports {
			exteriorPort := p.HostPort
			if exteriorPort == 0 {
				// No need to do port binding when HostPort is not specified
				continue
			}
			interiorPort := p.ContainerPort
			dockerPort := dockernat.Port(
				strconv.Itoa(int(interiorPort)) +
					toPortProtocol(p.Protocol))
			exposedPorts[dockerPort] = struct{}{}

			hostBinding := dockernat.PortBinding{
				HostPort: strconv.Itoa(int(exteriorPort)),
				HostIP:   p.HostIP,
			}
			// Allow multiple host ports bind to same docker port
			if existedBindings, ok := portBindings[dockerPort]; ok {
				// If a docker port already map to a host port, just append the host ports
				portBindings[dockerPort] = append(existedBindings, hostBinding)
			} else {
				// Otherwise, it's fresh new port binding
				portBindings[dockerPort] = []dockernat.PortBinding{
					hostBinding,
				}
			}
		}
	}

	sandboxConfig := &dockertypes.ContainerCreateConfig{
		Name: sandboxName,
		Config: &dockercontainer.Config{
			Image:        defaultPauseImage,
			ExposedPorts: exposedPorts,
			Labels:       newSandboxLabels(pod),
		},
		HostConfig: &dockercontainer.HostConfig{
			IpcMode:      dockercontainer.IpcMode("shareable"),
			PortBindings: portBindings,
		},
	}

	return sandboxConfig
}

func (m *cubeRuntimeManager) getAllPodsUID() ([]string, error) {
	filter := dockertypes.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", buildLabelSelector(ContainerTypeLabel, ContainerTypeSandbox)),
		),
	}

	sandboxes, err := m.dockerRuntime.ListContainers(filter)
	if err != nil {
		log.Printf("fail to list all pods sandbox: %v\n", err)
		return nil, err
	}

	// no build-in set in golang, using ugly code
	podUIDSet := make(map[string]bool)
	for _, sandbox := range sandboxes {
		if uid, ok := sandbox.Labels[PodUIDLabel]; ok {
			podUIDSet[uid] = true
		} else {
			log.Printf("[error] uid label of sandbox %s is empty\n", sandbox.Names[0])
		}
	}

	var uids []string
	for uid, _ := range podUIDSet {
		uids = append(uids, uid)
	}

	return uids, nil
}
