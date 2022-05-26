package types

type MQMessage struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	TriggerPath string `json:"triggerPath,omitempty" yaml:"triggerPath,omitempty"`

	ReturnTopic string `json:"returnTopic,omitempty" yaml:"returnTopic,omitempty"`
	ContentType string `json:"contentType" yaml:"contentType"`
	StatusCode  string `json:"statusCode" yaml:"statusCode"`

	Params  map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
	Payload string            `json:"payload" yaml:"payload"`
}
