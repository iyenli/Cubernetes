package proxyruntime

import (
	"Cubernetes/pkg/object"
	"fmt"
	"log"
	"net"

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
	FilterTable = "filter"
	NatTable    = "nat"
	InputChain  = "INPUT"
	OutputChain = "OUTPUT"
	DockerChain = "DOCKER"
	// SNAT use
	SnatOP      = "SNAT"
	PostRouting = "POSTROUTING"

	//DNAT use
	DnatOP     = "DNAT"
	PreRouting = "PREROUTING"
)

var ipt *iptables.IPTables

func InitIPTables() error {
	/* check env */
	var err error
	ipt, err = iptables.New(iptables.Timeout(3))
	if err != nil {
		log.Println(err)
		return err
	}

	var flag bool
	flag, err = ipt.ChainExists(FilterTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		return err
	}
	flag, err = ipt.ChainExists(NatTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		return err
	}
	return nil
}

func InitPodChain() error {
	return nil
}

func AddService(service *object.Service) error {
	// any service's cluster IP, modify to pod IP
	pods, err := GetPodByService(service)
	if err != nil {
		log.Println("Not matched pods found")
		return err
	}

	for _, pod := range pods {
		for _, port := range service.Spec.Ports {
			// push front as the highest priority
			// TODO: delete service and corresponding rules
			err = ipt.Insert(NatTable, PreRouting, 1,
				"-d", service.Spec.ClusterIP,
				"--dport", string(port.Port),
				"-p", string(port.Protocol),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", , string(port.TargetPort)))

			if err != nil {
				return err
			}
		}

	}
	return nil
}

// It would work even if the service not exist
func DeleteService(service *object.Service) error {
	pods, err := GetPodByService(service)
	if err != nil {
		log.Println("Not matched pods found")
		return err
	}

	for _, pod := range pods {
		for _, port := range service.Spec.Ports {
			// push front as the highest priority
			// TODO: delete service and corresponding rules
			err = ipt.Insert(NatTable, PreRouting, 1,
				"-d", service.Spec.ClusterIP,
				"--dport", string(port.Port),
				"-p", string(port.Protocol),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", pod.Status.IP.String(), string(port.TargetPort)))

			if err != nil {
				return err
			}
		}

	}
	return nil
}

func AddPod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			
			err := ipt.Insert(NatTable, PreRouting, 1,
				"-d", pod.Status.IP.String(),
				"--dport", string(port.HostPort),
				"-p", string(port.Protocol),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", dockerIP.String(), string(port.ContainerPort)))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//	ipTable.
//}
