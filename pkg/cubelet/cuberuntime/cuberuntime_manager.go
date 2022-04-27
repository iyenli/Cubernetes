package cuberuntime

import (
	"Cubernetes/pkg/cubelet/container"
	cubecontainer "Cubernetes/pkg/cubelet/container"
	dockershim "Cubernetes/pkg/cubelet/dockershim"
	object "Cubernetes/pkg/object"
	"log"
	"time"
)

const (
	containerdRuntimeName     = "containerd"
	podLogsRootDirectory      = "/var/log/pods"
	containerdRuntimeEndpoint = "unix:///run/containerd/containerd.sock"
	remoteConnectTimeout      = time.Second * 2
)

type cubeRuntimeManager struct {
	runtimeName string

	dockerRuntime dockershim.DockerRuntime
}

type podActions struct {
	KillPod           bool
	CreateSandbox     bool
	SandboxID         string
	ContainersToStart []int
	ContainersToKill  map[cubecontainer.ContainerID]*object.Container
}

func (m *cubeRuntimeManager) ListImages() ([]cubecontainer.Image, error) {
	var images []cubecontainer.Image

	allImages, err := m.dockerRuntime.ListImages(true)
	if err != nil {
		log.Printf("fail to list images\n")
		return images, err
	}

	for _, img := range allImages {
		images = append(images, cubecontainer.Image{
			ID:   img.ID,
			Size: int64(img.Size),
			Spec: cubecontainer.ImageSpec{
				Image: img.RepoTags[0],
			},
		})
	}

	return images, nil
}

func (m *cubeRuntimeManager) RemoveImage(image cubecontainer.ImageSpec) error {
	err := m.dockerRuntime.RemoveImage(image.Image)
	if err != nil {
		log.Printf("fail to remove image #{image.Name}\n")
		return err
	}

	return nil
}

type CubeRuntime interface {
	cubecontainer.Runtime
}

func (m *cubeRuntimeManager) SyncPod(pod *object.Pod, podStatus *cubecontainer.PodStatus) error {
	// Compute sandbox and container changes.
	podContainerChanges := m.computePodActions(pod, podStatus)

	removeContainer := false
	// Kill the pod if sandbox changed
	if podContainerChanges.KillPod {
		// kill all pod's containers
		m.killPodContainers(podStatus, removeContainer)

		// kill pod sandbox
		for _, sandbox := range podStatus.SandboxStatuses {
			if err := m.dockerRuntime.StopContainer(sandbox.Id); err != nil {
				log.Printf("fail to stop sandbox %s: %v\n", sandbox.Id, err)
				return err
			}

			if removeContainer {
				if err := m.dockerRuntime.RemoveContainer(sandbox.Id, false); err != nil {
					log.Printf("fail to remove sandbox %s: %v\n", sandbox.Id, err)
					return err
				}
			}
		}
	} else {
		// kill some containers
		for containerId, containerInfo := range podContainerChanges.ContainersToKill {
			if err := m.dockerRuntime.StopContainer(containerId.ID); err != nil {
				log.Printf("fail to kill container %s: %v\n", containerInfo.Name, err)
				return err
			}
		}
	}

	// Create sandbox if necessary
	podSandboxID := podContainerChanges.SandboxID
	podSandboxName := dockershim.MakeSandboxName(pod)
	if podContainerChanges.CreateSandbox {
		var err error

		if podSandboxName, podSandboxID, err = m.createPodSandbox(pod); err != nil {
			return err
		}
		log.Printf("create sandbox %s for pod %s\n", podSandboxID, pod.Name)
	}

	// Create containers
	for _, idx := range podContainerChanges.ContainersToStart {
		msg, err := m.startContainer(&pod.Spec.Containers[idx], pod, podSandboxName)
		if err != nil {
			log.Printf("fail to start container %s: %s\n", pod.Spec.Containers[idx].Name, msg)
			return err
		}
		log.Printf("start container %s\n", pod.Spec.Containers[idx].Name)
	}

	return nil
}

func (m *cubeRuntimeManager) computePodActions(pod *object.Pod, podStatus *cubecontainer.PodStatus) podActions {
	createPodSandbox, sandboxID := m.podSandboxChanged(pod, podStatus)
	changes := podActions{
		KillPod:           createPodSandbox,
		CreateSandbox:     createPodSandbox,
		SandboxID:         sandboxID,
		ContainersToStart: []int{},
		ContainersToKill:  make(map[cubecontainer.ContainerID]*object.Container),
	}

	// create sandbox need to (re-)create all containers
	if createPodSandbox {
		var containersToStart []int
		for idx := range pod.Spec.Containers {
			// TODO: RestartPolicy == OnFailure && ExitSucceeded => no need to start
			containersToStart = append(containersToStart, idx)
		}

		if len(containersToStart) == 0 {
			// nothing to create
			changes.CreateSandbox = false
			return changes
		}

		changes.ContainersToStart = containersToStart
		return changes
	}

	for idx, container := range pod.Spec.Containers {
		containerStatus := podStatus.FindContainerStatusByName(container.Name)

		if containerStatus == nil || containerStatus.State != cubecontainer.ContainerStateRunning {
			// container not exist or container not running
			if true /* TODO: container should be restart */ {
				changes.ContainersToStart = append(changes.ContainersToStart, idx)
				if containerStatus != nil && containerStatus.State == cubecontainer.ContainerStateUnknown {
					changes.ContainersToKill[containerStatus.ID] = &pod.Spec.Containers[idx]
				}
			}
		}

		// TODO: when we need to kill container?

	}

	return changes
}

// podSandboxChanged checks whether the spec of the pod is changed and returns
// (changed, original sandboxID if exist).
func (m *cubeRuntimeManager) podSandboxChanged(pod *object.Pod, podStatus *cubecontainer.PodStatus) (bool, string) {
	if len(podStatus.SandboxStatuses) == 0 {
		log.Printf("no sandbox for pod %s can be found. Need to start a new one.", pod.Name)
		// This branch should return
		return true, ""
	}

	sandboxStatus := podStatus.SandboxStatuses[0]
	if sandboxStatus.State != container.SandboxStateReady {
		// No ready sandbox for pod can be found. Need to start a new one.
		return true, sandboxStatus.Id
	}

	// Needs to create a new sandbox when the sandbox does not have an IP address.
	if sandboxStatus.Ip == "" {
		// Sandbox for pod has no IP address. Need to start a new one.
		return true, sandboxStatus.Id
	}

	return false, sandboxStatus.Id
}

func (m *cubeRuntimeManager) KillPod(UID string) error {
	podStatus, err := m.getPodStatusByUID(UID)
	if err != nil {
		log.Printf("fail to get podStatus by UID %s\n", UID)
		return err
	}
	// for debug only
	removeContainer := true

	// kill containers
	m.killPodContainers(podStatus, removeContainer)

	// kill pod sandbox
	for _, sandbox := range podStatus.SandboxStatuses {
		log.Printf("start to kill sandbox %s\n", sandbox.Id)
		if err := m.dockerRuntime.StopContainer(sandbox.Id); err != nil {
			log.Printf("fail to stop sandbox %s: %v\n", sandbox.Id, err)
			return err
		}

		if removeContainer {
			if err := m.dockerRuntime.RemoveContainer(sandbox.Id, false); err != nil {
				log.Printf("fail to remove sandbox %s: %v\n", sandbox.Id, err)
				return err
			}
		}
	}

	return nil
}

func (c *cubeRuntimeManager) getPodStatusByUID(UID string) (*container.PodStatus, error) {
	containerStatuses, err := c.getContainerStatusesByPodUID(UID)
	if err != nil {
		return nil, err
	}

	sandboxStatuses, err := c.getSandboxStatusesByPodUID(UID)
	if err != nil {
		return nil, err
	}

	podName := ""
	if len(sandboxStatuses) > 0 {
		podName = sandboxStatuses[0].Name
	}

	return &container.PodStatus{
		UID:               UID,
		Name:              podName,
		ContainerStatuses: containerStatuses,
		SandboxStatuses:   sandboxStatuses,
	}, nil
}

func NewCubeRuntimeManager() (CubeRuntime, error) {
	dockerRuntime, err := dockershim.NewDockerRuntime()
	if err != nil {
		log.Println("Fail to create docker client")
	}

	cm := &cubeRuntimeManager{
		dockerRuntime: dockerRuntime,
		runtimeName:   containerdRuntimeName,
	}

	return cm, nil
}

func (c *cubeRuntimeManager) Close() {
	c.dockerRuntime.CloseConnection()
}
