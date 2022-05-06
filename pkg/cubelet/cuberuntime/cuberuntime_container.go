package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/object"
	"log"
	"strings"
	"sync"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (m *cubeRuntimeManager) startContainer(container *object.Container, pod *object.Pod, podSandboxName string) (string, error) {
	err := m.dockerRuntime.PullImage(container.Image)
	if err != nil {
		log.Printf("ensure image for container #{container.Name} failed\n")
		return "", err
	}

	config := m.generateContainerConfig(container, pod, podSandboxName)
	containerID, err := m.dockerRuntime.CreateContainer(config)
	if err != nil {
		log.Printf("fail to create container #{container.Name}\n")
		return "", err
	}

	err = m.dockerRuntime.StartContainer(containerID)
	if err != nil {
		log.Printf("fail to start container #{container.Name}\n")
		return "", err
	}

	return containerID, nil
}

func (m *cubeRuntimeManager) getContainerStatusesByPodUID(UID string) ([]*cubecontainer.ContainerStatus, error) {
	filter := dockertypes.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", buildLabelSelector(ContainerTypeLabel, ContainerTypeContainer)),
			filters.Arg("label", buildLabelSelector(PodUIDLabel, UID)),
		),
	}

	containers, err := m.dockerRuntime.ListContainers(filter)
	if err != nil {
		log.Printf("fail to list pod containers %s: %v\n", UID, err)
		return nil, err
	}

	if len(containers) == 0 {
		return nil, nil
	}

	statuses := make([]*cubecontainer.ContainerStatus, len(containers))
	for i, container := range containers {
		statuses[i], _ = m.getContainerStatus(container.ID)
	}

	return statuses, nil
}

func (m *cubeRuntimeManager) generateContainerConfig(container *object.Container, pod *object.Pod, podSandboxName string) *dockertypes.ContainerCreateConfig {

	podContainerName := dockershim.MakeContainerName(pod, container)

	volumeBinds := make([]string, 0)
	for _, mount := range container.VolumeMounts {
		hostPath := findVolumeHostPath(pod, mount.Name)
		volumeBinds = append(volumeBinds,
			strings.Join([]string{hostPath, mount.MountPath}, ":"),
		)
	}

	mode := "container:" + podSandboxName

	config := &dockertypes.ContainerCreateConfig{
		Name: podContainerName,
		Config: &dockercontainer.Config{
			Image:  container.Image,
			Cmd:    container.Command,
			Labels: newContainerLabels(container, pod),
		},
		HostConfig: &dockercontainer.HostConfig{
			Binds:       volumeBinds,
			NetworkMode: dockercontainer.NetworkMode(mode),
			IpcMode:     dockercontainer.IpcMode(mode),
			PidMode:     dockercontainer.PidMode(mode),
		},
	}

	// set resource if specified
	if container.Resources != nil {
		config.HostConfig.Resources = dockercontainer.Resources{
			NanoCPUs: int64(container.Resources.Cpus * 1000000000),
			Memory:   container.Resources.Memory,
		}
	}

	return config
}

func (m *cubeRuntimeManager) killPodContainers(pod *cubecontainer.PodStatus, remove bool) {
	wg := sync.WaitGroup{}

	wg.Add(len(pod.ContainerStatuses))
	for _, container := range pod.ContainerStatuses {
		go func(container *cubecontainer.ContainerStatus) {
			defer wg.Done()

			if err := m.dockerRuntime.StopContainer(container.ID.ID); err != nil {
				log.Printf("error %v occurs when stoping container %s\n", err, container.Name)
			}

			if remove {
				if err := m.dockerRuntime.RemoveContainer(container.ID.ID, false); err != nil {
					log.Printf("error %v occurs when removing container %s\n", err, container.Name)
				}
			}
		}(container)
	}
	wg.Wait()
}

func (m *cubeRuntimeManager) getContainerStatus(UID string) (*cubecontainer.ContainerStatus, error) {
	containerJson, err := m.dockerRuntime.InspectContainer(UID)
	if err != nil {
		return nil, err
	}

	state, exitCode := toContainerStateAndExitCode(containerJson.State)
	created, _ := time.Parse(time.RFC3339Nano, containerJson.Created)
	started, _ := time.Parse(time.RFC3339Nano, containerJson.State.StartedAt)
	finished, _ := time.Parse(time.RFC3339Nano, containerJson.State.FinishedAt)

	statsJson, err := m.dockerRuntime.GetContainerStats(UID)
	if err != nil {
		return nil, err
	}

	status := &cubecontainer.ContainerStatus{
		ID: cubecontainer.ContainerID{
			Type: "docker",
			ID:   containerJson.ID,
		},
		Name:       dockershim.ParseContainerName(containerJson.Name),
		State:      state,
		CreatedAt:  created,
		StartedAt:  started,
		FinishedAt: finished,
		ResourceUsage: cubecontainer.ContainerResourceUsage{
			CPUUsage:    float64(statsJson.CPUStats.CPUUsage.TotalUsage) / 1000000000,
			MemoryUsage: int64(statsJson.MemoryStats.Usage),
		},
		ExitCode: exitCode,
		Image:    containerJson.Config.Image,
		ImageID:  strings.TrimLeft(containerJson.Image, "sha256:"),
	}

	return status, nil
}
