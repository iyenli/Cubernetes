package object

import (
	"net"
	"time"
)

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

type Ingress struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       IngressSpec    `json:"spec" yaml:"spec"`
	Status     *IngressStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type IngressSpec struct {
	// Http trigger used to call this Ingress
	TriggerPath string `json:"trigger" yaml:"trigger"`
	// Put json payload into this topic
	FeedChan string `json:"feed" yaml:"feed"`
	// Then wait json response from this topic
	ListenChan string `json:"listen" yaml:"listen"`
}

type IngressStatus struct {
	// nothing for now
}

type Actor struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ActionSpec   `json:"spec" yaml:"spec"`
	Status     *ActorStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ActorSpec struct {
	ActionName string `json:"action" yaml:"action"`
	ScriptFile string `json:"file" yaml:"file"`
}

type ActorStatus struct {
	Phase           ActorPhase `json:"phase,omitempty" yaml:"phase,omitempty"`
	IP              net.IP     `json:"IP" yaml:"IP"`
	NodeUID         string     `json:"node_uid,omitempty" yaml:"node_uid,omitempty"`
	LastUpdatedTime time.Time  `json:"lastUpdateTime" yaml:"lastUpdateTime"`
}

type ActorPhase string

const (
	ActorCreated ActorPhase = "created"
	ActorBound   ActorPhase = "bound"
	ActorRunning ActorPhase = "running"
	ActorFailed  ActorPhase = "failed"
	ActorUnknown ActorPhase = "unknown"
)
