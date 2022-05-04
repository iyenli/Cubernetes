package object

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

type PodStatus struct {
	// reserved for later use
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
	Memory int64 `json:"memory,omitempty" json:"memory,omitempty"`
}

type Volume struct {
	Name string `json:"name" yaml:"name"`
	// Volume only support HostPath type
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
	TargetPort int32    `json:"target,omitempty" yaml:"target,omitempty"`
}

type ServiceStatus struct {
	// reserved for later use
}
