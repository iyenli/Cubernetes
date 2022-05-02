package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubeproxy/utils"
	"Cubernetes/pkg/object"
	"errors"
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/google/uuid"
	"log"
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
	FilterTable  = "filter"
	NatTable     = "nat"
	InputChain   = "INPUT"
	OutputChain  = "OUTPUT"
	DockerChain  = "DOCKER"
	ServiceChain = "SERVICE"
	// SnatOP SNAT use
	SnatOP      = "SNAT"
	PostRouting = "POSTROUTING"

	// DnatOP DNAT use
	DnatOP     = "DNAT"
	PreRouting = "PREROUTING"

	// RANDOM Load balancer policy
	RANDOM    = "random"
	RR        = "nth"
	STATISTIC = "statistic"
)

var ipt *iptables.IPTables = nil
var serviceChainMap map[string]ServiceChainElement

// InitObject private function! Just for test
func InitObject() (err error) {
	ipt, err = iptables.New(iptables.Timeout(3))
	if err != nil {
		log.Println(err)
		return err
	}
	return
}

func InitIPTables() error {
	err := InitObject()
	if err != nil {
		return err
	}

	/* check env */
	flag, err := ipt.ChainExists(FilterTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		return err
	}
	flag, err = ipt.ChainExists(NatTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		return err
	}

	// Now, create SERVICE CHAIN, and add to PRE-ROUTING/OUTPUT Chain
	// Ref: https://gitee.com/k9-s/Cubernetes/wikis/IPT
	err = ipt.NewChain(NatTable, ServiceChain)
	if err != nil {
		log.Panicln("Creating chain failed")
		return err
	}

	err = ipt.Insert(NatTable, PreRouting,
		1, "-j", ServiceChain)
	if err != nil {
		log.Panicln("Creating chain failed")
		return err
	}

	err = ipt.Insert(NatTable, OutputChain, 1,
		"-j", ServiceChain)
	if err != nil {
		log.Panicln("Creating chain failed")
		return err
	}

	return nil
}

// ReleaseIPTables Delete all chains in service
func ReleaseIPTables() error {
	err := ipt.DeleteIfExists(NatTable, OutputChain, "-j", ServiceChain)
	if err != nil {
		log.Println("Error in release IP tables")
		return err
	}

	err = ipt.DeleteIfExists(NatTable, PreRouting, "-j", ServiceChain)
	if err != nil {
		log.Println("Error in release IP tables")
		return err
	}

	err = ipt.ClearAndDeleteChain(NatTable, ServiceChain)
	if err != nil {
		log.Println("Error in release IP tables")
		return err
	}
	return nil
}

func AddService(service *object.Service) error {
	// Default value of service
	// any service's cluster IP, modify to pod IP
	err := utils.DefaultService(service)
	if err != nil {
		return err
	}

	pods, err := GetPodByService(service)
	if err != nil || len(pods) == 0 {
		log.Println("Not matched pods found")
		return err
	}

	// init service chain element if NOT EXIST

	if _, ok := serviceChainMap[service.UID]; ok {
		err = DeleteService(service)
		if err != nil {
			log.Println("Delete service failed")
			return err
		}
	}

	prob := make([][]string, len(service.Spec.Ports))
	for idx, _ := range prob {
		prob[idx] = make([]string, len(pods))
	}

	serviceChainMap[service.UID] = ServiceChainElement{
		serviceChainUid:     []string{},
		probabilityChainUid: make([][]string, len(service.Spec.Ports)),
		numberOfPods:        len(pods),
	}

	for idx, port := range service.Spec.Ports {
		serviceUID := fmt.Sprintf("CUBE-SVC-%v", uuid.New().String())

		// Then create service chain and add to service
		err = ipt.NewChain(NatTable, serviceUID)
		if err != nil {
			log.Println("Create chain failed")
			return err
		}

		// depends on the settings
		// TODO: What if no port/protocol/target port specified?
		err := ipt.Insert(NatTable, ServiceChain,
			1, "-j", serviceUID,
			"-d", service.Spec.ClusterIP,
			"-p", string(port.Protocol),
			"--dport", strconv.FormatInt(int64(port.Port), 10))

		if err != nil {
			log.Panicln("Creating chain failed")
			return err
		}
		serviceChainMap[service.UID].serviceChainUid[idx] = serviceUID

		// Then create NUM(pod) chain
		for idx_, pod := range pods {
			podChainUID := fmt.Sprintf("CUBE-SVC-POD-%v", uuid.New().String())

			err = ipt.NewChain(NatTable, podChainUID)
			if err != nil {
				log.Println("Create chain failed")
				return err
			}

			// if 3 pods, the probability is 0.33/0.50/1.00, so...
			probability := float64(1) / float64(len(pods)-idx_)
			err = ipt.Insert(NatTable, serviceUID,
				1, "-j", podChainUID,
				"-m", STATISTIC,
				"--mode", RANDOM,
				"--probability", fmt.Sprintf("%.2f", probability),
			)
			if err != nil {
				log.Println("Create chain failed")
				return err
			}

			// at last, add DNAT service
			err = ipt.Insert(NatTable, podChainUID, 1,
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", pod.Status.IP.String(), strconv.FormatInt(int64(port.TargetPort), 10)),
			)

			if err != nil {
				log.Println("Create chain failed")
				return err
			}
			serviceChainMap[service.UID].probabilityChainUid[idx][idx_] = serviceUID

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
	if _, ok := serviceChainMap[service.UID]; !ok {
		log.Println("Delete not exist service")
		return errors.New("delete undef service")
	}

	// delete every
	for idx, port := range service.Spec.Ports {
		err := ipt.DeleteIfExists(NatTable, ServiceChain,
			"-j", serviceChainMap[service.UID].serviceChainUid[idx],
			"-d", service.Spec.ClusterIP,
			"-p", string(port.Protocol),
			"--dport", strconv.FormatInt(int64(port.Port), 10))

		if err != nil {
			log.Panicln("Deleting chain failed")
			return err
		}

		err = ipt.ClearAndDeleteChain(NatTable, serviceChainMap[service.UID].serviceChainUid[idx])
		if err != nil {
			log.Panicln("Deleting chain failed")
			return err
		}
	}

	for _, servicePort := range serviceChainMap[service.UID].probabilityChainUid {
		for _, dnat := range servicePort {
			err := ipt.ClearAndDeleteChain(NatTable, dnat)
			if err != nil {
				return err
			}
		}
	}

	// finally...
	delete(serviceChainMap, service.UID)
	return nil
}

func AddPodAsEndpoints(pod *object.Pod) error {
	services, err := crudobj.GetServices()
	if err != nil {
		log.Println("Get services failed")
		return err
	}

	for _, service := range services {
		if utils.MatchServiceAndPod(&service, pod) {
			// TODO: Finish this function
			err := reshuffleServiceIPTable(&service)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func reshuffleServiceIPTable(service *object.Service) error {
	return nil
}
