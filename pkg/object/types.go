package object

import (
	"net"
	"time"
)

type Pod struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       PodSpec `json:"spec" yaml:"spec"`
	// use pointer or else omitempty is disabled
	Status *PodStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type PodSpec struct {
	Containers []Container `json:"containers" yaml:"containers"`
	Volumes    []Volume    `json:"volumes,omitempty" yaml:"volumes,omitempty"`
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
	// reserved for later use
	IP                  net.IP         `json:"IP" yaml:"IP"`
	StartTime           time.Time      `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	Phase               PodPhase       `json:"phase,omitempty" yaml:"phase,omitempty"`
	NodeUID             string         `json:"pod-uid,omitempty" yaml:"pod-uid,omitempty"`
	ActualResourceUsage *ResourceUsage `json:"actualResourceUsage,omitempty" yaml:"actualResourceUsage,omitempty"`
	LastUpdateTime      time.Time      `json:"lastUpdateTime" yaml:"lastUpdateTime"`
}

type ResourceUsage struct {
	// for 4 cores, up to 400.00%
	ActualCPUUsage float64 `json:"actualCPUUsage" yaml:"actualCPUUsage"`
	// in bytes
	ActualMemoryUsage int64 `json:"actualMemoryUsage" yaml:"actualMemoryUsage"`
}

type Container struct {
	Name    string   `json:"name" yaml:"name"`
	Image   string   `json:"image" yaml:"image"`
	Command []string `json:"command,omitempty" yaml:"command,omitempty"`
	Args    []string `json:"args,omitempty" yaml:"args,omitempty"`
	// use pointer or else omitempty is disabled
	Resources    *ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
	VolumeMounts []VolumeMount         `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
	Ports        []ContainerPort       `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type ResourceRequirements struct {
	Cpus float64 `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	// Memory in bytes
	Memory int64 `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type Volume struct {
	Name string `json:"name" yaml:"name"`
	// Volume only support HostPath types
	HostPath string `json:"hostPath" yaml:"hostPath"`
}

type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
}

type ContainerPort struct {
	Name          string `json:"name" yaml:"name"`
	HostPort      int32  `json:"hostPort" yaml:"hostPort"`
	ContainerPort int32  `json:"containerPort" yaml:"containerPort"`
	Protocol      string `json:"protocol" yaml:"protocol"`
	HostIP        string `json:"hostIP" yaml:"hostIP"`
}

type Service struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ServiceSpec    `json:"spec" yaml:"spec"`
	Status     *ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ServiceSpec struct {
	Selector  map[string]string `json:"selector,omitempty" yaml:"selector,omitempty"`
	Ports     []ServicePort     `json:"ports,omitempty" yaml:"ports,omitempty"`
	ClusterIP string            `json:"ip,omitempty" yaml:"ip,omitempty"`
}

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

type ServicePort struct {
	Protocol   Protocol `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Port       int32    `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort int32    `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
}

type ServiceStatus struct {
	Endpoints []net.IP     `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Ingress   []PodIngress `json:"ingress,omitempty" yaml:"ingress,omitempty"`
}

type PodIngress struct {
	HostName string  `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IP       string  `json:"ip,omitempty" yaml:"ip,omitempty"`
	Ports    []int32 `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type ReplicaSet struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ReplicaSetSpec    `json:"spec" yaml:"spec"`
	Status     *ReplicaSetStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ReplicaSetSpec struct {
	Replicas int32             `json:"replicas" yaml:"replicas"`
	Selector map[string]string `json:"selector,omitempty" yaml:"selector,omitempty"`
	Template PodTemplate       `json:"template" yaml:"template"`
}

type PodTemplate struct {
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       PodSpec `json:"spec" yaml:"spec"`
}

type ReplicaSetStatus struct {
	// actual running pod replica in PodUIDs
	RunningReplicas int32 `json:"replicas" yaml:"replicas"`
	// UID of pods assigned
	PodUIDsToRun   []string  `json:"podsToRun" yaml:"podsToRun"`
	PodUIDsToKill  []string  `json:"podsToKill" yaml:"podsTokill"`
	PodUIDsRunning []string  `json:"pods" yaml:"pods"`
	LastUpdateTime time.Time `json:"lastUpdate,omitempty" yaml:"lastUpdate,omitempty"`
}

type Node struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       NodeSpec    `json:"spec" yaml:"spec"`
	Status     *NodeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type NodeType string

const (
	Master NodeType = "Master"
	Slave  NodeType = "Slave"
)

type NodeSpec struct {
	Type     NodeType     `json:"types" yaml:"types"`
	Capacity NodeCapacity `json:"capacity,omitempty" yaml:"capacity,omitempty"`
	Info     NodeInfo     `json:"info,omitempty" yaml:"info,omitempty"`
}

type NodeStatus struct {
	Addresses NodeAddresses `json:"addresses,omitempty" yaml:"addresses,omitempty"`
	Condition NodeCondition `json:"condition,omitempty" yaml:"condition,omitempty"`
}

type NodeAddresses struct {
	HostName   string `json:"hostName,omitempty" yaml:"hostName,omitempty"`
	ExternalIP string `json:"externalIP,omitempty" yaml:"externalIP,omitempty"`
	InternalIP string `json:"internalIP,omitempty" yaml:"internalIP,omitempty"`
}

type NodeCondition struct {
	OutOfDisk      bool `json:"outOfDisk" yaml:"outOfDisk"`
	Ready          bool `json:"ready" yaml:"ready"`
	MemoryPressure bool `json:"memoryPressure" yaml:"memoryPressure"`
	DiskPressure   bool `json:"diskPressure" yaml:"diskPressure"`
}

type NodeCapacity struct {
	CPUCount int `json:"cpuCount,omitempty" yaml:"cpuCount,omitempty"`
	Memory   int `json:"memory,omitempty" yaml:"memory,omitempty"`
	MaxPods  int `json:"maxPods,omitempty" yaml:"maxPods,omitempty"`
}

type NodeInfo struct {
	CubeVersion   string `json:"cubeVersion,omitempty" yaml:"cubeVersion,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty" yaml:"kernelVersion,omitempty"`
	DeviceName    string `json:"deviceName,omitempty" yaml:"deviceName,omitempty"`
}

type AutoScaler struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
}

type AutoScalerSpec struct {
	// Target audience for scale, should always be Pod
	Workload string `json:"workload" yaml:"workload"`
	// lower limit for the number of pods that can be set by the autoscaler, default 1.
	MinReplicas int `json:"minReplicas" yaml:"minReplicas"`

	// upper limit for the number of pods that can be set by the autoscaler.
	// It cannot be smaller than MinReplicas.
	MaxReplicas int `json:"maxReplicas" yaml:"maxReplicas"`

	// target average CPU utilization & Memory bytes over all the pods.
	// nil if not specified.
	TargetUtilization Utilization `json:"targetUtilization" yaml:"targetUtilization"`
}

type AutoScalerStatus struct {
	LastScaleTime     time.Time   `json:"lastScale,omitempty" yaml:"lastScale,omitempty"`
	LastUpdateTime    time.Time   `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	ReplicaSetUID     string      `json:"replicaSetUID" yaml:"replicaSetUID"`
	DesiredReplicas   int         `json:"desiredReplicas" yaml:"desiredReplicas"`
	ActualReplicas    int         `json:"actualReplicas" yaml:"actualReplicas"`
	ActualUtilization Utilization `json:"actualUtilization" yaml:"actualUtilization"`
}

type Utilization struct {
	CPU    *CpuUtilization    `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory *MemoryUtilization `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type CpuUtilization struct {
	MinPercentage float64 `json:"minPercentage" yaml:"minPercentage"`
	MaxPercentage float64 `json:"maxPercentage" yaml:"maxPercentage"`
}

type MemoryUtilization struct {
	MinBytes int64 `json:"minBytes" yaml:"minBytes"`
	MaxBytes int64 `json:"maxBytes" yaml:"maxBytes"`
}
type Dns struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       DnsSpec `json:"spec" yaml:"spec"`
}

type DnsSpec struct {
	Host  string                    `json:"host" yaml:"host"`
	Paths map[string]DnsDestination `json:"paths" yaml:"paths"`
}

type DnsDestination struct {
	ServiceUID  string `json:"serviceUID" yaml:"serviceUID"`
	ServicePort int32  `json:"servicePort" yaml:"servicePort"`
}
