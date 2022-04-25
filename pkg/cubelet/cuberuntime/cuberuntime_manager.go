package cuberuntime

import (
	cubecontainer "Cubernetes/pkg/cubelet/container"
	images "Cubernetes/pkg/cubelet/images"
	object "Cubernetes/pkg/object"
	"log"

	criapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const (
	containerdRuntimeName = "containerd"
	podLogsRootDirectory  = "/var/log/pods"
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
	CreateSandbox     bool
	SandboxID         string
	Attempt           uint32
	ContainersToStart []int
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
}



func (m *cubeRuntimeManager) SyncPod(pod *object.Pod, podStatus *cubecontainer.PodStatus) (string, error) {
	// Compute sandbox and container changes.
	podContainerChanges := m.computePodActions(pod, podStatus)

	// Create sandbox if necessary
	podSandboxID := podContainerChanges.SandboxID
	if podContainerChanges.CreateSandbox {
		var msg string
		var err error

		podSandboxID, msg, err = m.createPodSandbox(pod, podContainerChanges.Attempt)
		if err != nil {
			return msg, err
		}
		log.Printf("create sandbox %s for pod %s\n", podSandboxID, pod.Name)
	}

	// TODO: create containers, kill pod, kill containers...

	return "success", nil
}

func (m *cubeRuntimeManager) computePodActions(pod *object.Pod, podStatus *cubecontainer.PodStatus) podActions {
	createPodSandbox, attempt, sandboxID := m.podSandboxChanged(pod, podStatus)
	changes := podActions{
		CreateSandbox:     createPodSandbox,
		SandboxID:         sandboxID,
		Attempt:           attempt,
		ContainersToStart: []int{},
	}

	// TODO: calculate containers to start

	return changes
}

// podSandboxChanged checks whether the spec of the pod is changed and returns
// (changed, new attempt, original sandboxID if exist).
func (m *cubeRuntimeManager) podSandboxChanged(pod *object.Pod, podStatus *cubecontainer.PodStatus) (bool, uint32, string) {
	if podStatus.SandboxStatus == nil {
		// No sandbox for pod can be found. Need to start a new one.
		// This branch should return 
		return true, 0, ""
	}

	sandboxStatus := podStatus.SandboxStatus
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
