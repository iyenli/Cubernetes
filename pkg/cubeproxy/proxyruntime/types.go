package proxyruntime

import (
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubeproxy/informer"
	"github.com/coreos/go-iptables/iptables"
)

/**
@Chenfan
							IPTables
                               XXXXXXXXXXXXXXXXXX
                             XXX     Network    XXX
                               XXXXXXXXXXXXXXXXXX
                                       +
                                       |
                                       v
 +-------------+              +------------------+
 |table: filter| <---+        | table: nat       |
 |chain: INPUT |     |        | chain: PREROUTING|
 +-----+-------+     |        +--------+---------+
       |             |                 |
       v             |                 v
 [local process]     |           ****************          +--------------+
       |             +---------+ Routing decision +------> |table: filter |
       v                         ****************          |chain: FORWARD|
****************                                           +------+-------+
Routing decision                                                  |
****************                                                  |
       |                                                          |
       v                        ****************                  |
+-------------+       +------>  Routing decision  <---------------+
|table: nat   |       |         ****************
|chain: OUTPUT|       |               +
+-----+-------+       |               |
      |               |               v
      v               |      +-------------------+
+--------------+      |      | table: nat        |
|table: filter | +----+      | chain: POSTROUTING|
|chain: OUTPUT |             +--------+----------+
+--------------+                      |
                                      v
                               XXXXXXXXXXXXXXXXXX
                             XXX    Network     XXX
                               XXXXXXXXXXXXXXXXXX
*/

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
	ProbabilityChainUid [][]string
	ServiceChainUid     []string
	NumberOfPods        int

	RelevantDNS []string
}

type DNSElement struct {
	ContainerID string
}

type ProxyRuntime struct {
	Ipt            *iptables.IPTables
	DockerInstance dockershim.DockerRuntime

	PodInformer     informer.PodInformer
	ServiceInformer informer.ServiceInformer
	DNSInformer     informer.DNSInformer

	ServiceChainMap map[string]ServiceChainElement
	DNSMap          map[string]DNSElement
}
