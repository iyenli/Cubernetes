package object

type MQMessage struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	ReturnTopic string `json:"returnTopic" yaml:"returnTopic"`
	TriggerPath string `json:"triggerPath" yaml:"triggerPath"`

	Body string `json:"body,omitempty" yaml:"body,omitempty"`
}
