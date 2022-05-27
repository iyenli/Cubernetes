package types

type MQMessage struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	TriggerPath string `json:"triggerPath" yaml:"triggerPath"`

	ReturnTopic string `json:"returnTopic" yaml:"returnTopic"`
	ContentType string `json:"contentType" yaml:"contentType"`
	StatusCode  string `json:"statusCode" yaml:"statusCode"`

	Params  map[string]string `json:"params" yaml:"params"`
	Payload string            `json:"payload" yaml:"payload"`
}
