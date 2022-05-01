package object

type TypeMeta struct {
	Kind       string `json,yaml:"kind,omitempty"`
	APIVersion string `json,yaml:"apiVersion,omitempty"`
}

type ObjectMeta struct {
	Name      string            `json,yaml:"name"`
	Namespace string            `json,yaml:"namespace,omitempty"`
	UID       string            `json,yaml:"uid,omitempty"`
	Labels    map[string]string `json,yaml:"labels,omitempty"`
}
