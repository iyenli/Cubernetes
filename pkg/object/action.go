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
	// Action names that this Action would invoke
	InvokeActions []string `json:"invokeActions" yaml:"invokeActions"`
}

type ActionStatus struct {
	Phase           ActionPhase `json:"phase" yaml:"phase"`
	LastScaleTime   time.Time   `json:"lastScale,omitempty" yaml:"lastScale,omitempty"`
	LastUpdateTime  time.Time   `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	DesiredReplicas int         `json:"desiredReplicas" yaml:"desiredReplicas"`
	ActualReplicas  int         `json:"actualReplicas" yaml:"actualReplicas"`

	Actors []string `json:"actors" yaml:"actors"`
	ToRun  []string `json:"toRun" yaml:"toRun"`
	ToKill []string `json:"toKill" yaml:"toKill"`
}

type ActionPhase string

const (
	ActionCreating ActionPhase = "Creating"
	ActionCreated  ActionPhase = "Created"
	ActionServing  ActionPhase = "Serving"
)
