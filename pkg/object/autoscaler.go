package object

import "time"

const AutoScalerEtcdPrefix = "/apis/autoScaler/"

type AutoScaler struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       AutoScalerSpec    `json:"spec" yaml:"spec"`
	Status     *AutoScalerStatus `json:"status,omitempty" yaml:"status,omitempty"`
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
