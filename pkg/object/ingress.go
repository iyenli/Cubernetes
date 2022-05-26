package object

const IngressEtcdPrefix = "/apis/ingress/"

type IngressPhase string

const (
	IngressCreating IngressPhase = "Creating"
	IngressReady    IngressPhase = "Ready"
)

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
	InvokeAction string `json:"invokeAction" yaml:"invokeAction"`
	// http request type
	HTTPType string `json:"httpType,omitempty" yaml:"httpType,omitempty"`
}

type IngressStatus struct {
	// nothing for now
	Phase IngressPhase `yaml:"phase,omitempty" json:"phase,omitempty"`
}
