package object

import "time"

const ActionEtcdPrefix = "/apis/action/"

type Action struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ActionSpec    `json:"spec" yaml:"spec"`
	Status     *ActionStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ActionSpec struct {
	// Host path of corresponding python script
	// manager should also copy the script to /etc/cubernetes/actions
	// so that user can't modify script while action is running
	ScriptPath string `json:"scriptPath" yaml:"scriptPath"`
	// Topic name that this Action subscribe
	InputChan string `json:"input" yaml:"input"`
	// Topic name that this Action output to
	OutputChan string `json:"output" yaml:"output"`
}

type ActionStatus struct {
	// Implement by ReplicaSet, kinda like AutoScaler
	ReplicaSetUID   string    `json:"replicaSetUID" yaml:"replicaSetUID"`
	LastScaleTime   time.Time `json:"lastScale,omitempty" yaml:"lastScale,omitempty"`
	LastUpdateTime  time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	DesiredReplicas int       `json:"desiredReplicas" yaml:"desiredReplicas"`
	ActualReplicas  int       `json:"actualReplicas" yaml:"actualReplicas"`
}
