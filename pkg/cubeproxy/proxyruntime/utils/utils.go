package utils

import (
	"Cubernetes/pkg/object"
	"errors"
	"log"
	"net"
)

var (
	DNSConfigError = errors.New("illegal hostname")
)

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

func CheckDNS(dns *object.Dns) error {
	host := dns.Spec.Host

	if len(host) == 0 {
		return DNSConfigError
	}

	// Hostname be like: xxx.xxx.xxx
	if host[0] == '/' {
		host = host[1:]
	}
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	portMap := make(map[int32]bool)
	pathMap := make(map[string]bool)
	for path, dst := range dns.Spec.Paths {
		if _, ok := pathMap[path]; ok {
			return DNSConfigError
		}
		if _, ok := portMap[dst.ServicePort]; ok {
			return DNSConfigError
		}
		portMap[dst.ServicePort] = true
		pathMap[path] = true
	}

	return nil
}
