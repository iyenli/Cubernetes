package cuberuntime

import (
	"Cubernetes/pkg/object"
	"log"
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
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

	return podSandboxConfig.Name, sandboxID, nil

}

func generatePodSandboxConfig(pod *object.Pod) *dockertypes.ContainerCreateConfig {
	sandboxName := pod.Name + "_sandbox"

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
			// only support tcp now
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
		},
		HostConfig: &dockercontainer.HostConfig{
			IpcMode:      dockercontainer.IpcMode("shareable"),
			PortBindings: portBindings,
		},
	}

	return sandboxConfig
}

