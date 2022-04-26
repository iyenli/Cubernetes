package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/object"
	"log"
	"sync"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	killContainerTimeout = 2
)

func (m *cubeRuntimeManager) startContainer(podSandboxID string, podSandboxConfig *runtimeapi.PodSandboxConfig,
	container *object.Container, pod *object.Pod, podStatus *cubecontainer.PodStatus, podIP string, podIPs []string) (string, error) {
	imageRef, err := m.imagePuller.EnsureImageExists(pod, container, podSandboxConfig)
	if err != nil {
		log.Printf("ensure image for container #{container.Name} failed\n")
		return "", err
	}

	config, _ := m.generateContainerConfig(container, pod, podIP, imageRef, podIPs)
	containerID, err := m.runtimeService.CreateContainer(podSandboxID, config, podSandboxConfig)
	if err != nil {
		log.Printf("fail to create container #{container.Name}\n")
		return "", err
	}

	return containerID, nil
}

func (m *cubeRuntimeManager) generateContainerConfig(container *object.Container, pod *object.Pod, podIP, imageRef string, podIPs []string) (*runtimeapi.ContainerConfig, error) {

	// this generated config is NOT compelete, not even close,
	// I just filled what make sense to me.
	config := &runtimeapi.ContainerConfig{
		Metadata: &runtimeapi.ContainerMetadata{
			Name:    container.Name,
			Attempt: 0, // 0 before I figure out its usage
		},
		Image:   &runtimeapi.ImageSpec{Image: imageRef},
		Command: container.Command,
		Args:    container.Args,
		Labels:  newContainerLabels(container, pod),
		Mounts: generateContainerMounts(pod, container),
	}

	return config, nil
}

func (m *cubeRuntimeManager) killPodContainers(runningPod cubecontainer.Pod) {
	wg := sync.WaitGroup{}

	wg.Add(len(runningPod.Containers))
	for _, container := range runningPod.Containers {
		go func(container *cubecontainer.Container) {
			defer wg.Done()

			err := m.runtimeService.StopContainer(container.ID.ID, killContainerTimeout)
			if err != nil {
				log.Printf("error %v occurs when killing container %s\n", err, container.Name)
			}
		}(container)
	}
	wg.Wait()
}
