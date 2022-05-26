package object

type MQMessage struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	TriggerPath string `json:"triggerPath" yaml:"triggerPath"`

	ReturnTopic string `json:"returnTopic" yaml:"returnTopic"`
	ReturnType  string `json:"ReturnType" yaml:"ReturnType"`

	Params map[string]string `json:"Params" yaml:"Params"`
	Body   string            `json:"body,omitempty" yaml:"body,omitempty"`
}
