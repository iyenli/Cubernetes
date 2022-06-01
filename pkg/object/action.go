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
	ScriptUID string `json:"scriptUID" yaml:"scriptUID"`
	// Action names that this Action would invoke
	InvokeActions []string `json:"invokeActions" yaml:"invokeActions"`
}

type ActionStatus struct {
	LastScaleTime   time.Time `json:"lastScale,omitempty" yaml:"lastScale,omitempty"`
	LastUpdateTime  time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	DesiredReplicas int       `json:"desiredReplicas" yaml:"desiredReplicas"`
	ActualReplicas  int       `json:"actualReplicas" yaml:"actualReplicas"`

	Actors []string `json:"actors" yaml:"actors"`
	ToRun  []string `json:"toRun" yaml:"toRun"`
	ToKill []string `json:"toKill" yaml:"toKill"`
}
