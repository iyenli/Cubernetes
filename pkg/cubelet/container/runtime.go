package container

import (
	object "Cubernetes/pkg/object"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"time"
)

// Runtime interface defines the interfaces that should be implemented
// by a container runtime.
type Runtime interface {
	Type() string
	GetPods() ([]*Pod, error)
	GetPodStatus(uid, name, namespace string)
	SyncPod(pod *object.Pod, podStatus *PodStatus)
}

type ContainerID struct {
	Type string
	ID   string
}

type Container struct {
	ID      ContainerID
	Name    string
	Image   string
	ImageID string
	Hash    uint64
	State   string
}

type ContainerStatus struct {
	ID         ContainerID
	Name       string
	State      string
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	ExitCode   int
	Image      string
	ImageID    string
	Hash       uint64
	Reason     string
	Message    string
}

type Pod struct {
	UID        string
	Name       string
	Namespace  string
	Containers []*Container
	SandBoxes  []*Container
}

type PodStatus struct {
	UID               string
	Name              string
	Namespace         string
	IPs               []string
	ContainerStatuses []*ContainerStatus
	SandboxStatuses   []*runtimeapi.PodSandboxStatus
}

// Annotation represents an annotation.
type Annotation struct {
	Name  string
	Value string
}

// ImageSpec describes a specified image with annotations.
type ImageSpec struct {
	Image       string
	Annotations []Annotation
}

type Image struct {
	ID   string
	Size int64
	Spec ImageSpec
}

type ImageService interface {
	PullImage(image ImageSpec, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error)
	GetImageRef(image ImageSpec) (string, error)
	ListImages() ([]Image, error)
	RemoveImage(image ImageSpec) error
}
