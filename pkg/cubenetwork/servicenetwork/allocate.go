package servicenetwork

import (
	"Cubernetes/pkg/cubenetwork/servicenetwork/utils"
	"Cubernetes/pkg/object"
	"log"
	"net"
)

func (cia *ClusterIPAllocator) AllocateClusterIP(service *object.Service) (object.Service, error) {
	ip := net.ParseIP(service.Spec.ClusterIP)
	if ip != nil && cia.ipNet.Contains(ip) {
		// entry not exist or labeled removed
		if exist, ok := cia.mp[utils.Ip2int(ip)]; !ok || !exist {
			cia.mp[utils.Ip2int(ip)] = true
			return *service, nil
		}
	}

	// else: overwrite it
	hasAllocated := false
	toAllocate := cia.nextIP
	for !hasAllocated {
		cia.nextIP++
		if !cia.ipNet.Contains(utils.Int2ip(toAllocate)) {
			log.Fatal("Cubernetes use all ips for service")
		}

		if exist, ok := cia.mp[toAllocate]; !ok || !exist {
			cia.mp[toAllocate] = true
			hasAllocated = true
		}
	}

	service.Spec.ClusterIP = utils.Int2ip(toAllocate).String()
	log.Printf("Allocate IP for service, service uid is %v, ip is %v", service.UID, service.Spec.ClusterIP)

	return *service, nil
}

func (cia *ClusterIPAllocator) DeAllocateClusterIP(service *object.Service) error {
	cia.mp[utils.Ip2int(net.ParseIP(service.Spec.ClusterIP))] = false
	return nil
}
