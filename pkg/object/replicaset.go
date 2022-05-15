package object

import "time"

const ReplicaSetEtcdPrefix = "/apis/replicaSet/"

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
