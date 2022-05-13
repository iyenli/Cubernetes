package proxyruntime

import (
	"Cubernetes/pkg/object"
	"errors"
	"github.com/coreos/go-iptables/iptables"
	"log"
	"net"
)

func (pr *ProxyRuntime) ClearAllService() error {
	for _, service := range pr.ServiceInformer.ListServices() {
		err := pr.DeleteService(&service)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitObject private function! Just for test
func (pr *ProxyRuntime) InitObject() (err error) {
	pr.Ipt, err = iptables.New(iptables.Timeout(3))
	if err != nil {
		log.Println(err)
		return err
	}
	return
}

// CheckService
// 1. Check if service is legal
// 2. Make it legal using default values in K8
func CheckService(service *object.Service) error {
	if net.ParseIP(service.Spec.ClusterIP) == nil {
		log.Println("[Fatal]: Illegal Cluster IP")
		return errors.New("illegal cluster ip")
	}

	for _, port := range service.Spec.Ports {
		if port.Port == 0 {
			port.Port = port.TargetPort
		}
		if port.Protocol == "" {
			port.Protocol = object.ProtocolTCP
		}
	}

	return nil
}
