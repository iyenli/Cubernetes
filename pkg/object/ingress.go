package object

const IngressEtcdPrefix = "/apis/ingress/"

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
