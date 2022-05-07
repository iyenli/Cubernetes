package container

import (
	object "Cubernetes/pkg/object"
	"net"
	"time"
)

// Runtime interface defines the interfaces that should be implemented
// by a container runtime.
type Runtime interface {
	// GetPodStatus SyncPod Type() string
	// GetPods() ([]*Pod, error)
	GetPodStatus(UID string) (*PodStatus, error)
	KillPod(UID string) error
	SyncPod(pod *object.Pod, podStatus *PodStatus) error
	InspectPod(UID string) (*object.PodStatus, error)
	ListPodsUID() ([]string, error)

	Close()
}

type ContainerID struct {
	Type string
	ID   string
}

type ContainerState string
type SandboxState string

const (
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
	ID            ContainerID
	Name          string
	State         ContainerState
	CreatedAt     time.Time
	StartedAt     time.Time
	FinishedAt    time.Time
	ResourceUsage ContainerResourceUsage
	ExitCode      int
	Image         string
	ImageID       string
}

type ContainerResourceUsage struct {
	CPUUsage    float64
	MemoryUsage int64
}

type SandboxStatus struct {
	Id     string
	Name   string
	PodUID string
	State  SandboxState
	Ip     string
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

func (s *PodStatus) FindContainerStatusByName(containerName string) *ContainerStatus {
	for _, containerStatus := range s.ContainerStatuses {
		if containerStatus.Name == containerName {
			return containerStatus
		}
	}
	return nil
}

func (s *PodStatus) UpdateSandboxStatuses(sandboxStatuses []*SandboxStatus) {
	s.SandboxStatuses = sandboxStatuses
}

func ComputePodPhase(statuses []*ContainerStatus, sandboxStatus *SandboxStatus) object.PodPhase {
	isRunning := true
	runningContainer := 0
	isSuccess := true
	isFail := false

	for _, container := range statuses {
		if container.State == ContainerStateRunning {
			runningContainer += 1
		} else if container.State != ContainerStateCreated {
			isRunning = false
		}
		if container.State != ContainerStateExited || container.ExitCode != 0 {
			isSuccess = false
			if container.State == ContainerStateExited {
				isFail = true
			}
		}
	}

	if isFail {
		return object.PodFailed
	} else if isSuccess {
		return object.PodSucceeded
	} else if isRunning && sandboxStatus.State == SandboxStateReady {
		return object.PodRunning
	} else if runningContainer > 0 && sandboxStatus.State == SandboxStateReady {
		return object.PodPending
	} else {
		return object.PodUnknown
	}
}
