package object

const (
	KindPod        = "Pod"
	KindService    = "Service"
	KindReplicaSet = "ReplicaSet"
	KindNode       = "Node"
	KindDns        = "Dns"
	KindAutoScaler = "AutoScaler"
	KindGpuJob     = "GpuJob"
)

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

type ObjectMeta struct {
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	UID         string            `json:"uid,omitempty" yaml:"uid,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
