package utils

import (
	"Cubernetes/pkg/object"
	"errors"
	"log"
	"net"
)

var (
	HostnameError = errors.New("illegal hostname")
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

func CheckDNSHostName(host string) (string, error) {
	if len(host) == 0 {
		return "", HostnameError
	}

	if host[0] == '/' {
		host = host[1:]
	}
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}
	return host, nil
}
