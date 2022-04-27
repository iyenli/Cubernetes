package container

import (
	object "Cubernetes/pkg/object"
	"net"
	"time"
)

// Runtime interface defines the interfaces that should be implemented
// by a container runtime.
type Runtime interface {
	// SyncPod Type() string
	// GetPods() ([]*Pod, error)
	// GetPodStatus(uid, name, namespace string) (*PodStatus, error)
	KillPod(UID string) error
	SyncPod(pod *object.Pod, podStatus *PodStatus) error

	Close()
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

type ContainerState string
type SandboxState string

const (
	runtimeName = "docker"
	// ContainerStateCreated indicates a container that has been created (e.g. with docker create) but not started.
	ContainerStateCreated ContainerState = "created"
	// ContainerStateRunning indicates a currently running container.
	ContainerStateRunning ContainerState = "running"
	// ContainerStateExited indicates a container that ran and completed ("stopped" in other contexts, although a created container is technically also "stopped").
	ContainerStateExited ContainerState = "exited"
	// ContainerStateUnknown encompasses all the states that we currently don't care about (like restarting, paused, dead).
	ContainerStateUnknown ContainerState = "unknown"

	SandboxStateReady    SandboxState = "ready"
	SandboxStateNotReady SandboxState = "not ready"
)

type ContainerStatus struct {
	ID         ContainerID
	Name       string
	State      ContainerState
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

type SandboxStatus struct {
	Id     string
	Name   string
	PodUID string
	State  SandboxState
	Ip     string
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
	NetworkNamespace  string
	PodNetWork        PodNetworkStatus
	ContainerStatuses []*ContainerStatus
	SandboxStatuses   []*SandboxStatus
}

// PodNetworkStatus stores the network status of a pod (currently just the primary IP address)
// This struct represents version "v1beta1"
type PodNetworkStatus struct {
	object.TypeMeta `json:",inline"`

	// IP is the primary ipv4/ipv6 address of the pod. Among other things it is the address that -
	//   - kube expects to be reachable across the cluster
	//   - service endpoints are constructed with
	//   - will be reported in the PodStatus.PodIP field (will override the IP reported by docker)
	IP net.IP `json:"ip" description:"Primary IP address of the pod"`
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

func (podStatus *PodStatus) FindContainerStatusByName(containerName string) *ContainerStatus {
	for _, containerStatus := range podStatus.ContainerStatuses {
		if containerStatus.Name == containerName {
			return containerStatus
		}
	}
	return nil
}

func ConvertPodStatusToRunningPod(podStatus *PodStatus) Pod {
	runningPod := Pod{
		UID:       podStatus.UID,
		Name:      podStatus.Name,
		Namespace: podStatus.Namespace,
	}

	for _, containerStatus := range podStatus.ContainerStatuses {
		if containerStatus.State != ContainerStateRunning {
			continue
		}
		container := &Container{
			ID:      containerStatus.ID,
			Name:    containerStatus.Name,
			Image:   containerStatus.Image,
			ImageID: containerStatus.ImageID,
			Hash:    containerStatus.Hash,
		}
		runningPod.Containers = append(runningPod.Containers, container)
	}

	for _, sandbox := range podStatus.SandboxStatuses {
		runningPod.SandBoxes = append(runningPod.SandBoxes, &Container{
			ID:    ContainerID{Type: runtimeName, ID: sandbox.Id},
			State: string(sandbox.State),
		})
	}

	return runningPod
}
