package servicenetwork

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubenetwork/servicenetwork/utils"
	"log"
	"net"
)

type ClusterIPAllocator struct {
	mp map[uint32]bool

	nextIP uint32
	ipNet  net.IPNet
}

func NewClusterIPAllocator() *ClusterIPAllocator {
	cia := ClusterIPAllocator{
		mp:     map[uint32]bool{},
		nextIP: 0,
		ipNet:  net.IPNet{},
	}
	err := cia.Init()
	if err != nil {
		log.Println("Allocate cluster ip allocator failed")
		return nil
	}

	return &cia
}

func (cia *ClusterIPAllocator) Init() error {
	nextIP, ipNet, err := net.ParseCIDR(cubeconfig.ServiceClusterIPRange)
	if err != nil {
		log.Println("Init cluster ip allocator failed")
		return err
	}

	cia.nextIP = utils.Ip2int(nextIP)
	cia.ipNet = *ipNet

	services, err := crudobj.GetServices()
	if err != nil {
		log.Println("Get service failed when allocate cluster ip allocator")
		return err
	}

	for _, service := range services {
		ip := net.ParseIP(service.Spec.ClusterIP)
		if ip == nil {
			log.Println("[INFO] Service exist but don't have cluster IP")
			continue
		}

		if !cia.ipNet.Contains(ip) {
			log.Fatal("[Fatal]: Exist service IP not in CIDR")
		}

		IPNumber := utils.Ip2int(ip)
		cia.mp[IPNumber] = true
	}

	return nil
}
