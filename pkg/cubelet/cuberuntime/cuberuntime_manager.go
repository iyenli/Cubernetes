package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	cri "Cubernetes/pkg/cubelet/cri"
	images "Cubernetes/pkg/cubelet/images"
	object "Cubernetes/pkg/object"
	"log"
	"time"

	criapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	containerdRuntimeName     = "containerd"
	podLogsRootDirectory      = "/var/log/pods"
	containerdRuntimeEndpoint = "unix:///run/containerd/containerd.sock"
	remoteConnectTimeout      = time.Second * 2
)

type cubeRuntimeManager struct {
	runtimeName string

	// wrapped image puller.
	imagePuller images.ImageManager

	// grpc service client
	runtimeService criapi.RuntimeService
	imageService   criapi.ImageManagerService
}

type podActions struct {
	KillPod           bool
	CreateSandbox     bool
	SandboxID         string
	Attempt           uint32
	ContainersToStart []int
	ContainersToKill  map[cubecontainer.ContainerID]*object.Container
}

func (m *cubeRuntimeManager) PullImage(image cubecontainer.ImageSpec, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	// Pull without AuthConfig: not supported
	imageRef, err := m.imageService.PullImage(toRuntimeAPIImageSpec(image), nil, podSandboxConfig)
	if err != nil {
		log.Printf("fail to pull image #{image.Name}\n")
		return "", err
	}

	return imageRef, nil
}

func (m *cubeRuntimeManager) GetImageRef(image cubecontainer.ImageSpec) (string, error) {
	status, err := m.imageService.ImageStatus(toRuntimeAPIImageSpec(image))
	if err != nil {
		log.Printf("fail to get image #{image.Name} status\n")
		return "", err
	}

	if status == nil {
		return "", nil
	}

	return status.Id, nil
}

func (m *cubeRuntimeManager) ListImages() ([]cubecontainer.Image, error) {
	var images []cubecontainer.Image

	allImages, err := m.imageService.ListImages(nil)
	if err != nil {
		log.Printf("fail to list images\n")
		return images, err
	}

	for _, img := range allImages {
		images = append(images, cubecontainer.Image{
			ID:   img.Id,
			Size: int64(img.Size_),
			Spec: toCubeContainerImageSpec(img),
		})
	}

	return images, nil
}

func (m *cubeRuntimeManager) RemoveImage(image cubecontainer.ImageSpec) error {
	err := m.imageService.RemoveImage(&runtimeapi.ImageSpec{Image: image.Image})
	if err != nil {
		log.Printf("fail to remove image #{image.Name}\n")
		return err
	}

	return nil
}

type CubeRuntime interface {
	cubecontainer.Runtime
	cubecontainer.ImageService
}

func (m *cubeRuntimeManager) SyncPod(pod *object.Pod, podStatus *cubecontainer.PodStatus) error {
	// Compute sandbox and container changes.
	podContainerChanges := m.computePodActions(pod, podStatus)

	// Kill the pod if sandbox changed
	if podContainerChanges.KillPod {
		// kill all pod's containers
		m.killPodContainers(cubecontainer.ConvertPodStatusToRunningPod(podStatus))

		// kill pod sandbox
		for _, sandbox := range podStatus.SandboxStatuses {
			if err := m.runtimeService.StopPodSandbox(sandbox.Id); err != nil {
				log.Printf("fail to kill sandbox %s: %v\n", sandbox.Id, err)
				return err
			}
		}
	} else {
		// kill some containers
		for containerId, containerInfo := range podContainerChanges.ContainersToKill {
			if err := m.runtimeService.StopContainer(containerId.ID, killContainerTimeout); err != nil {
				log.Printf("fail to kill container %s: %v\n", containerInfo.Name, err)
				return err
			}
		}
	}

	var podIPs []string
	if podStatus != nil {
		podIPs = podStatus.IPs
	}
	// Create sandbox if necessary
	podSandboxID := podContainerChanges.SandboxID
	if podContainerChanges.CreateSandbox {
		var err error

		if podSandboxID, _, err = m.createPodSandbox(pod, podContainerChanges.Attempt); err != nil {
			return err
		}
		log.Printf("create sandbox %s for pod %s\n", podSandboxID, pod.Name)
		podSandboxStatus, err := m.runtimeService.PodSandboxStatus(podSandboxID)
		if err != nil {
			log.Printf("fail to get pod status %s: %v\n", podSandboxID, err)
			return err
		}
		podIPs = m.determinePodSandboxIPs(pod.Namespace, pod.Name, podSandboxStatus)
	}

	podIP := ""
	if len(podIPs) != 0 {
		podIP = podIPs[0]
	}

	podSandboxConfig, _ := m.generatePodSandboxConfig(pod, podContainerChanges.Attempt)
	// Create containers
	for _, idx := range podContainerChanges.ContainersToStart {
		msg, err := m.startContainer(podSandboxID, podSandboxConfig, &pod.Spec.Containers[idx], pod, podStatus, podIP, podIPs)
		if err != nil {
			log.Printf("fail to start container %s: %s\n", pod.Spec.Containers[idx].Name, msg)
			return err
		}
		log.Printf("start container %s\n", pod.Spec.Containers[idx].Name)
	}

	return nil
}

func (m *cubeRuntimeManager) computePodActions(pod *object.Pod, podStatus *cubecontainer.PodStatus) podActions {
	createPodSandbox, attempt, sandboxID := m.podSandboxChanged(pod, podStatus)
	changes := podActions{
		KillPod:           createPodSandbox,
		CreateSandbox:     createPodSandbox,
		SandboxID:         sandboxID,
		Attempt:           attempt,
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
// (changed, new attempt, original sandboxID if exist).
func (m *cubeRuntimeManager) podSandboxChanged(pod *object.Pod, podStatus *cubecontainer.PodStatus) (bool, uint32, string) {
	if len(podStatus.SandboxStatuses) == 0 {
		// No sandbox for pod can be found. Need to start a new one.
		// This branch should return
		return true, 0, ""
	}

	sandboxStatus := podStatus.SandboxStatuses[0]
	if sandboxStatus.State != runtimeapi.PodSandboxState_SANDBOX_READY {
		// No ready sandbox for pod can be found. Need to start a new one.
		return true, sandboxStatus.Metadata.Attempt + 1, sandboxStatus.Id
	}

	// Needs to create a new sandbox when network namespace changed.
	if sandboxStatus.GetLinux().GetNamespaces().GetOptions().GetNetwork() != runtimeapi.NamespaceMode_POD {
		// Sandbox for pod has changed. Need to start a new one.
		return true, sandboxStatus.Metadata.Attempt + 1, ""
	}

	// Needs to create a new sandbox when the sandbox does not have an IP address.
	if sandboxStatus.Network.Ip == "" {
		// Sandbox for pod has no IP address. Need to start a new one.
		return true, sandboxStatus.Metadata.Attempt + 1, sandboxStatus.Id
	}

	return false, sandboxStatus.Metadata.Attempt, sandboxStatus.Id
}

// func (m *cubeRuntimeManager) GetPodStatus(uid, name, namespace string) (*cubecontainer.PodStatus, error) {

// }

func NewCubeRuntimeManager() (CubeRuntime, error) {
	runtimeService, err := cri.NewRemoteRuntimeService(containerdRuntimeEndpoint, remoteConnectTimeout)
	if err != nil {
		return nil, err
	}

	imageService, err := cri.NewRemoteImageService(containerdRuntimeEndpoint, remoteConnectTimeout)
	if err != nil {
		return nil, err
	}

	cm := &cubeRuntimeManager{
		runtimeName:    containerdRuntimeName,
		runtimeService: runtimeService,
		imageService:   imageService,
	}

	cm.imagePuller = images.NewImageManager(cm)

	return cm, nil
}
