package proxyruntime

import (
	"github.com/coreos/go-iptables/iptables"
	"log"
)

// CleanIptables Just an util used in cuberoot
func CleanIptables() error {
	ipt, err := iptables.New(iptables.Timeout(3))

	chains, err := ipt.ListChains("nat")
	if err != nil {
		log.Println("[Error]: list chains failed")
		return err
	}

	err = ipt.ClearChain("nat", ServiceChain)
	if err != nil {
		log.Println("[Error]: clear service chain failed")
		return err
	}

	// For all SVC Chains
	for _, chain := range chains {
		if len(chain) >= len("CUBE-SVC") && chain[:8] == "CUBE-SVC" {
			err := ipt.ClearChain(NatTable, chain)
			if err != nil {
				log.Println("[Error]: clear cube-service chain failed")
				return err
			}
		}
	}

	for _, chain := range chains {
		if len(chain) >= len("CUBE-SVC") && chain[:8] == "CUBE-SVC" {
			err := ipt.DeleteChain(NatTable, chain)
			if err != nil {
				log.Println("[Error]: delete cube-service chain failed")
				return err
			}
		}
	}

	// delete redundant service chain
	for {
		if exist, err := ipt.Exists(NatTable, PreRouting, "-j", ServiceChain); err != nil || !exist {
			break
		}
		err = ipt.Delete(NatTable, PreRouting, "-j", ServiceChain)
		if err != nil {
			return err
		}
	}
	for {
		if exist, err := ipt.Exists(NatTable, OutputChain, "-j", ServiceChain); err != nil || !exist {
			break
		}
		err := ipt.Delete(NatTable, OutputChain, "-j", ServiceChain)
		if err != nil {
			return err
		}
	}

	return nil
}
