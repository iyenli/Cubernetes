package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubeproxy/informer"
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

func (pr *ProxyRuntime) MapPortToPods(service *object.Service, pods []object.Pod, port *object.ServicePort, idx int) error {
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
	for idx_, pod := range pods {
		podChainUID := fmt.Sprintf("CUBE-SVC-POD-%v", uuid.New().String()[:15])

		err = pr.Ipt.NewChain(NatTable, podChainUID)
		if err != nil {
			log.Println("Create pod probability chain failed")
			return err
		}

		// if 3 pods, the probability is 0.33/0.50/1.00, so...
		if idx_ < len(pods)-1 {
			probability := float64(1) / float64(len(pods)-idx_)
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
			"--to-destination", fmt.Sprintf("%v:%v", pod.Status.IP.String(),
				strconv.FormatInt(int64(port.TargetPort), 10)),
		)

		if err != nil {
			log.Println("Create chain failed")
			return err
		}
		pr.ServiceChainMap[service.UID].probabilityChainUid[idx][idx_] = podChainUID

		err = crudobj.AddEndpointToService(service, pod.Status.IP)
		if err != nil {
			log.Println("Update endpoint IP to API Server failed")
			return err
		}
	}
	return nil
}

func InitProxyRuntime() (*ProxyRuntime, error) {
	pr := &ProxyRuntime{
		Ipt:             nil,
		ServiceChainMap: make(map[string]ServiceChainElement),
		ServiceInformer: informer.NewServiceInformer(),
		PodInformer:     informer.NewPodInformer(),
	}

	err := pr.InitObject()
	if err != nil {
		log.Panicln("Init object failed")
		return nil, err
	}

	/* check env */
	flag, err := pr.Ipt.ChainExists(FilterTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		//return nil, err
	}
	flag, err = pr.Ipt.ChainExists(NatTable, DockerChain)
	if !flag {
		log.Printf("Start docker first")
		//return nil, err
	}
	/* Check env ends */

	// Clear all service chain:
	for exist, err := pr.Ipt.Exists(NatTable, PreRouting, "-j", ServiceChain); err != nil && exist; {
		err := pr.Ipt.Delete(NatTable, PreRouting, "-j", ServiceChain)
		if err != nil {
			return nil, err
		}
	}
	for exist, err := pr.Ipt.Exists(NatTable, OutputChain, "-j", ServiceChain); err != nil && exist; {
		err := pr.Ipt.Delete(NatTable, OutputChain, "-j", ServiceChain)
		if err != nil {
			return nil, err
		}
	}

	// create SERVICE CHAIN, and add to PRE-ROUTING/OUTPUT Chain
	// Ref: https://gitee.com/k9-s/Cubernetes/wikis/IPT
	if exists, _ := pr.Ipt.ChainExists(NatTable, ServiceChain); !exists {
		err = pr.Ipt.NewChain(NatTable, ServiceChain)
		if err != nil {
			log.Panicln("[Panic]: Creating chain failed")
			return nil, err
		}
	}

	err = pr.Ipt.Insert(NatTable, PreRouting,
		1, "-j", ServiceChain)
	if err != nil {
		log.Panicln("[Panic]: Add chain failed")
		return nil, err
	}

	err = pr.Ipt.Insert(NatTable, OutputChain, 1,
		"-j", ServiceChain)
	if err != nil {
		log.Panicln("[Panic]: Add chain failed")
		return nil, err
	}

	return pr, nil
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
