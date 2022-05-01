package object

type Pod struct {
	TypeMeta   `json,yaml:",inline"`
	ObjectMeta `json,yaml:"metadata"`
	Spec       PodSpec   `json,yaml:"spec"`
	Status     PodStatus `json,yaml:"status,omitempty"`
}

type PodSpec struct {
	Containers []Container `json,yaml:"containers"`
	Volumes    []Volume    `json,yaml:"volumes,omitempty"`
}

type PodStatus struct {
	// reserved for later use
}

type Container struct {
	Name         string               `json,yaml:"name"`
	Image        string               `json,yaml:"image"`
	Command      []string             `json,yaml:"command,omitempty"`
	Resources    ResourceRequirements `json,yaml:"resources,omitempty"`
	VolumeMounts []VolumeMount        `json,yaml:"volumeMounts,omitempty"`
	Ports        []ContainerPort      `json,yaml:"ports,omitempty"`
}

type ResourceRequirements struct {
	Limits   map[string]string `json,yaml:"limits,omitempty"`
	Requests map[string]string `json,yaml:"requests,omitempty"`
}

type Volume struct {
	Name string `json,yaml:"name"`
	// Volume only support HostPath type
	HostPath string `json,yaml:"hostPath"`
}

type VolumeMount struct {
	Name      string `json,yaml:"name"`
	MountPath string `json,yaml:"mountPath"`
}

type ContainerPort struct {
	Name          string `json,yaml:"name"`
	HostPort      int32  `json,yaml:"hostPort"`
	ContainerPort int32  `json,yaml:"containerPort"`
	Protocol      string `json,yaml:"protocol"`
}
