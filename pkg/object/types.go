package object

type Pod struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       PodSpec `json:"spec"`
	// use pointer or else omitempty is disabled
	Status *PodStatus `json:"status,omitempty"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
	Volumes    []Volume    `json:"volumes,omitempty"`
}

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodCreated means that API Pod object was created by API Server
	PodCreated PodPhase = "Created"
	// PodBound means that scheduler bind this pod to a node
	PodBound PodPhase = "Bound"
	// PodAccepted means cubelet has accepted this pod to run
	PodAccepted PodPhase = "Accepted"
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

type PodStatus struct {
	Phase PodPhase `json:"phase,omitempty"`
}

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	// use pointer or else omitempty is disabled
	Resources    *ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts []VolumeMount         `json:"volumeMounts,omitempty"`
	Ports        []ContainerPort       `json:"ports,omitempty"`
}

type ResourceRequirements struct {
	Cpus float64 `json:"cpus,omitempty"`
	// Memory in bytes
	Memory int64 `json:"memory,omitempty"`
}

type Volume struct {
	Name string `json:"name"`
	// Volume only support HostPath type
	HostPath string `json:"hostPath"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type ContainerPort struct {
	Name          string `json:"name"`
	HostPort      int32  `json:"hostPort"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP"`
}

type Service struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       ServiceSpec   `json:"spec"`
	Status     ServiceStatus `json:"status,omitempty"`
}

type ServiceSpec struct {
	Selector  map[string]string `json:"selector,omitempty"`
	Ports     []ServicePort     `json:"ports,omitempty"`
	ClusterIP string            `json:"ip,omitempty"`
}

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

type ServicePort struct {
	Protocol   Protocol `json:"protocol,omitempty"`
	Port       int32    `json:"port,omitempty"`
	TargetPort int32    `json:"target,omitempty"`
}

type ServiceStatus struct {
	Ingress []PodIngress `json:"ingress,omitempty"`
}

type PodIngress struct {
	HostName string  `json:"hostname,omitempty"`
	IP       string  `json:"ip,omitempty"`
	Ports    []int32 `json:"ports,omitempty"`
}

type ReplicaSet struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       ReplicaSetSpec    `json:"spec"`
	Status     *ReplicaSetStatus `json:"status,omitempty"`
}

type ReplicaSetSpec struct {
	Replicas int32             `json:"replicas"`
	Selector map[string]string `json:"selector,omitempty"`
	Template PodTemplate       `json:"template"`
}

type PodTemplate struct {
	ObjectMeta `json:"metadata"`
	Spec       PodSpec `json:"spec"`
}

type ReplicaSetStatus struct {
	// actual runnig pod replica in PodUIDs
	RunningReplicas int32 `json:"replicas"`
	// UID of pods assigned
	PodUIDs []string `json:"podUIDs"`
}
