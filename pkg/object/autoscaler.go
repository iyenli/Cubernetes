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

	Template PodTemplate `json:"template" yaml:"template"`
	// lower limit for the number of pods that can be set by the autoscaler, default 1.
	MinReplicas int `json:"minReplicas" yaml:"minReplicas"`

	// upper limit for the number of pods that can be set by the autoscaler.
	// It cannot be smaller than MinReplicas.
	MaxReplicas int `json:"maxReplicas" yaml:"maxReplicas"`

	// target average CPU utilization & Memory bytes over all the pods.
	// nil if not specified.
	TargetUtilization UtilizationLimit `json:"targetUtilization" yaml:"targetUtilization"`
}

type AutoScalerStatus struct {
	LastScaleTime     time.Time          `json:"lastScale,omitempty" yaml:"lastScale,omitempty"`
	LastUpdateTime    time.Time          `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	ReplicaSetUID     string             `json:"replicaSetUID" yaml:"replicaSetUID"`
	DesiredReplicas   int                `json:"desiredReplicas" yaml:"desiredReplicas"`
	ActualReplicas    int                `json:"actualReplicas" yaml:"actualReplicas"`
	ActualUtilization AverageUtilization `json:"actualUtilization" yaml:"actualUtilization"`
}

type UtilizationLimit struct {
	CPU    *CpuUtilizationLimit    `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory *MemoryUtilizationLimit `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type CpuUtilizationLimit struct {
	MinPercentage float64 `json:"minPercentage" yaml:"minPercentage"`
	MaxPercentage float64 `json:"maxPercentage" yaml:"maxPercentage"`
}

type MemoryUtilizationLimit struct {
	MinBytes int64 `json:"minBytes" yaml:"minBytes"`
	MaxBytes int64 `json:"maxBytes" yaml:"maxBytes"`
}

type AverageUtilization struct {
	CPUPercentage float64 `json:"cpuPercentage" yaml:"cpuPercentage"`
	MemoryBytes   int64   `json:"memoryBytes" yaml:"memoryBytes"`
}
