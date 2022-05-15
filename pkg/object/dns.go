package object

type Dns struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       DnsSpec `json:"spec" yaml:"spec"`
}

type DnsSpec struct {
	Host  string                    `json:"host" yaml:"host"`
	Paths map[string]DnsDestination `json:"paths" yaml:"paths"`
}

type DnsDestination struct {
	ServiceUID  string `json:"serviceUID" yaml:"serviceUID"`
	ServicePort int32  `json:"servicePort" yaml:"servicePort"`
}
