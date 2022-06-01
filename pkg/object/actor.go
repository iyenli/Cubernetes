package object

import (
	"net"
	"time"
)

const ActorEtcdPrefix = "/apis/actor/"

type Actor struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ActorSpec    `json:"spec" yaml:"spec"`
	Status     *ActorStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ActorSpec struct {
	ActionName    string   `json:"actionName" yaml:"actionName"`
	ScriptUID     string   `json:"file" yaml:"file"`
	InvokeActions []string `json:"invokeActions" yaml:"invokeActions"`
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
