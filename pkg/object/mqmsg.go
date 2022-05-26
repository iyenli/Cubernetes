package object

type MQMessage struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	TriggerPath string `json:"triggerPath,omitempty" yaml:"triggerPath,omitempty"`

	ReturnTopic string `json:"returnTopic,omitempty" yaml:"returnTopic,omitempty"`
	ReturnType  string `json:"returnType,omitempty" yaml:"returnType,omitempty"`

	Params map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
	Body   string            `json:"body,omitempty" yaml:"body,omitempty"`
}
