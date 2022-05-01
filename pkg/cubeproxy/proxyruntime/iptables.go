package proxyruntime

import (
	"Cubernetes/pkg/object"
	"github.com/coreos/go-iptables/iptables"
	"log"
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

func AddService(service *object.Service) {
	ipTable, err := iptables.New(iptables.Timeout(3))
	if err != nil {
		return
	}

	_, err = ipTable.ChainExists(NatTable, DockerChain)
	if err != nil {
		return
	}
}

//	ipTable.
//}
