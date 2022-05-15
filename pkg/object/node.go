package object

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
