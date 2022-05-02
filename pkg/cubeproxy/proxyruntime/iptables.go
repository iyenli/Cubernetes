package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"log"
	"net"
	"strconv"
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
	// SnatOP SNAT use
	SnatOP      = "SNAT"
	PostRouting = "POSTROUTING"

	// DnatOP DNAT use
	DnatOP     = "DNAT"
	PreRouting = "PREROUTING"

	// RANDOM Load balancer policy
	RANDOM = "random"
	RR     = "nth"
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
	// Default value of service
	// any service's cluster IP, modify to pod IP
	pods, err := GetPodByService(service)
	if err != nil || len(pods) == 0 {
		log.Println("Not matched pods found")
		return err
	}

	// Easy load balancer
	podNums := len(pods)
	probability := float64(1) / float64(podNums)

	// Load balancer: Random and average
	for _, pod := range pods {
		for _, port := range service.Spec.Ports {
			// push front as the highest priority
			// TODO: delete service and corresponding rules
			err = ipt.Insert(NatTable, PreRouting, 1,
				"-d", service.Spec.ClusterIP,
				"-p", string(port.Protocol),
				"--dport", strconv.FormatInt(int64(port.TargetPort), 10),
				"--mode", RANDOM,
				"--probability", fmt.Sprintf("%.2f", probability),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", pod.Status.IP.String(), strconv.FormatInt(int64(port.Port), 10)))

			if err != nil {
				return err
			}
			err = crudobj.AddEndpointToService(service, pod.Status.IP)
			if err != nil {
				log.Println("Update endpoint IP to API Server failed")
				return err
			}
		}

	}
	return nil
}

// DeleteService It would work even if the service not exist
func DeleteService(service *object.Service) error {
	pods, err := GetPodByService(service)
	if err != nil || len(pods) == 0 {
		log.Println("Not matched pods found")
		return err
	}

	for _, pod := range pods {
		for _, port := range service.Spec.Ports {
			// push front as the highest priority
			// TODO: delete service and corresponding rules
			err = ipt.DeleteIfExists(NatTable, PreRouting,
				"-d", service.Spec.ClusterIP,
				"-p", string(port.Protocol),
				"--dport", strconv.FormatInt(int64(port.Port), 10),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", pod.Status.IP.String(), strconv.FormatInt(int64(port.TargetPort), 10)))

			if err != nil {
				return err
			}
		}

	}
	return nil
}

// AddPod FIX: -p should set front of --dport
func AddPod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			err := ipt.Append(NatTable, PreRouting,
				"-d", pod.Status.IP.String(),
				"-p", port.Protocol,
				"--dport", strconv.FormatInt(int64(port.HostPort), 10),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", dockerIP.String(), strconv.FormatInt(int64(port.ContainerPort), 10)))

			if err != nil {
				log.Println("Add pod IP to iptables failed")
				return err
			}
		}
	}

	return nil
}

// DeletePod FIX: -p should set front of --dport
func DeletePod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {

			err := ipt.DeleteIfExists(NatTable, PreRouting,
				"-d", pod.Status.IP.String(),
				"-p", port.Protocol,
				"--dport", strconv.FormatInt(int64(port.HostPort), 10),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", dockerIP.String(), strconv.FormatInt(int64(port.ContainerPort), 10)))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//	ipTable.
//}
