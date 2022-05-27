package message

type MonitorMessage struct {
	InvokeTimeUnix int64  `json:"time" yaml:"time"`
	Action         string `json:"action" yaml:"action"`
}
