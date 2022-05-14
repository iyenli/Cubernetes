package proxyruntime

import (
	"Cubernetes/pkg/object"
	"fmt"
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

func (pr *ProxyRuntime) MapPortToPods(service *object.Service, podIPs []string, port *object.ServicePort, idx int) error {
	// Chain name under 29 chars
	serviceUID := fmt.Sprintf("CUBE-SVC-%v", uuid.New().String()[:15])

	// Then create service chain and add to service
	err := pr.Ipt.NewChain(NatTable, serviceUID)
	if err != nil {
		log.Println("Create service chain failed")
		return err
	}

	err = pr.Ipt.Insert(NatTable, ServiceChain,
		1, "-j", serviceUID,
		"-d", service.Spec.ClusterIP,
		"-p", string(port.Protocol),
		"--dport", strconv.FormatInt(int64(port.Port), 10))
	if err != nil {
		log.Panicln("Add service chain to service failed")
		return err
	}

	pr.ServiceChainMap[service.UID].serviceChainUid[idx] = serviceUID

	// Then create NUM(pod) chain
	for idx_, pod := range podIPs {
		podChainUID := fmt.Sprintf("CUBE-SVC-POD-%v", uuid.New().String()[:15])

		err = pr.Ipt.NewChain(NatTable, podChainUID)
		if err != nil {
			log.Println("Create pod probability chain failed")
			return err
		}

		// if 3 podIPs, the probability is 0.33/0.50/1.00, so...
		if idx_ < len(podIPs)-1 {
			probability := float64(1) / float64(len(podIPs)-idx_)
			err = pr.Ipt.Append(NatTable, serviceUID,
				"-j", podChainUID,
				"-m", STATISTIC,
				"--mode", RANDOM,
				"--probability", fmt.Sprintf("%.2f", probability),
			)
			if err != nil {
				log.Println("Add probability chain to service chain failed")
				return err
			}
		} else {
			err = pr.Ipt.Append(NatTable, serviceUID,
				"-j", podChainUID,
			)
			if err != nil {
				log.Println("Add probability chain to service chain failed")
				return err
			}
		}

		// at last, add DNAT service
		err = pr.Ipt.Insert(NatTable, podChainUID, 1,
			"-j", DnatOP,
			"-p", string(port.Protocol),
			"--to-destination", fmt.Sprintf("%v:%v", pod,
				strconv.FormatInt(int64(port.TargetPort), 10)),
		)

		if err != nil {
			log.Println("Create chain failed")
			return err
		}

		pr.ServiceChainMap[service.UID].probabilityChainUid[idx][idx_] = podChainUID
	}
	return nil
}

// ReleaseIPTables Delete all chains in service
func (pr *ProxyRuntime) ReleaseIPTables() error {
	if !TestPurpose {
		err := pr.ClearAllService()
		if err != nil {
			log.Panicln("Error in clear all service.")
			return err
		}
	}

	if exists, _ := pr.Ipt.ChainExists(NatTable, ServiceChain); exists {
		err := pr.Ipt.Delete(NatTable, OutputChain, "-j", ServiceChain)

		if err != nil {
			log.Println("Error in release IP tables")
			return err
		}

		err = pr.Ipt.Delete(NatTable, PreRouting, "-j", ServiceChain)
		if err != nil {
			log.Println("Error in release IP tables")
			return err
		}

		err = pr.Ipt.ClearAndDeleteChain(NatTable, ServiceChain)
		if err != nil {
			log.Println("Error in release IP tables")
			return err
		}
	}
	return nil
}
