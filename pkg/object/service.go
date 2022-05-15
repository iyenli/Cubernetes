package object

import "net"

const ServiceEtcdPrefix = "/apis/service/"

type Service struct {
	TypeMeta   `json:",inline" yaml:",inline"`
	ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       ServiceSpec    `json:"spec" yaml:"spec"`
	Status     *ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ServiceSpec struct {
	Selector  map[string]string `json:"selector,omitempty" yaml:"selector,omitempty"`
	Ports     []ServicePort     `json:"ports,omitempty" yaml:"ports,omitempty"`
	ClusterIP string            `json:"ip,omitempty" yaml:"ip,omitempty"`
}

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

type ServicePort struct {
	Protocol   Protocol `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Port       int32    `json:"port,omitempty" yaml:"port,omitempty"`
	TargetPort int32    `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
}

type ServiceStatus struct {
	Endpoints []net.IP     `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Ingress   []PodIngress `json:"ingress,omitempty" yaml:"ingress,omitempty"`
}

type PodIngress struct {
	HostName string  `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IP       string  `json:"ip,omitempty" yaml:"ip,omitempty"`
	Ports    []int32 `json:"ports,omitempty" yaml:"ports,omitempty"`
}
