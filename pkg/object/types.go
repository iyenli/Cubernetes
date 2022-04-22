package object

type Pod struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
	Spec       PodSpec   `json:"spec"`
	Status     PodStatus `json:"status,omitempty"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
	Volumes    []Volume    `json:"volumes,omitempty"`
}

type PodStatus struct {
	// reserved for later use
}

type Container struct {
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Command      []string             `json:"command,omitempty"`
	Args         []string             `json:"args,omitempty"`
	Resources    ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts []VolumeMount        `json:"volumeMounts,omitempty"`
	Ports        []ContainerPort      `json:"ports,omitempty"`
}

type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
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
}
