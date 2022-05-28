package types

type MQMessageRequest struct {
	RequestUID  string `json:"requestUID" yaml:"requestUID"`
	TriggerPath string `json:"triggerPath" yaml:"triggerPath"`
	ReturnTopic string `json:"returnTopic" yaml:"returnTopic"`

	Params  map[string]string `json:"params" yaml:"params"`
	Payload string            `json:"payload" yaml:"payload"`
}

type MQMessageResponse struct {
	RequestUID string `json:"requestUID" yaml:"requestUID"`

	ContentType string `json:"contentType" yaml:"contentType"`
	StatusCode  string `json:"statusCode" yaml:"statusCode"`

	Payload string `json:"payload" yaml:"payload"`
}
