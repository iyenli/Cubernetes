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
	InvokeAction string `json:"feed" yaml:"feed"`
	// http request type
	HTTPType string `json:"HTTPType" yaml:"HTTPType"`
}

type IngressStatus struct {
	// nothing for now
	Phase IngressPhase `yaml:"phase" json:"phase"`
}
