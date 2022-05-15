package proxyruntime

import (
	"Cubernetes/pkg/cubeproxy/informer"
	"github.com/coreos/go-iptables/iptables"
)

const (
	FilterTable  = "filter"
	NatTable     = "nat"
	InputChain   = "INPUT"
	OutputChain  = "OUTPUT"
	DockerChain  = "DOCKER"
	ServiceChain = "SERVICE"

	// DnatOP DNAT use
	DnatOP     = "DNAT"
	PreRouting = "PREROUTING"

	// RANDOM Load balancer policy
	RANDOM      = "random"
	STATISTIC   = "statistic"
	TestPurpose = false
)

type ServiceChainElement struct {
	probabilityChainUid [][]string
	serviceChainUid     []string
	numberOfPods        int
}

type DNSElement struct {
}

type ProxyRuntime struct {
	Ipt             *iptables.IPTables
	PodInformer     informer.PodInformer
	ServiceInformer informer.ServiceInformer
	DNSInformer     informer.DNSInformer

	ServiceChainMap map[string]ServiceChainElement
	DNSMap          map[string]DNSElement
}
