package proxyruntime

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime/utils"
	"Cubernetes/pkg/object"
	"errors"
	"log"
	"strconv"
)

func (pr *ProxyRuntime) AddDNS(dns *object.Dns) error {
	// preprocess hostname, e.g. example.com is ok
	// /example.com.xx/ not ok
	err := utils.CheckDNS(dns)
	if err != nil {
		log.Fatalln("[Fatal]: DNS host config is not legal")
		return err
	}
	hostname := dns.Spec.Host

	paths := make([]string, len(dns.Spec.Paths))
	serviceIP := make([]string, len(dns.Spec.Paths))
	port := make([]string, len(dns.Spec.Paths))
	services := pr.ServiceInformer.ListServices()

	index := 0
	for path, dst := range dns.Spec.Paths {
		paths[index] = path
		port[index] = strconv.FormatInt(int64(dst.ServicePort), 10)

		serviceExist := false
		for _, service := range services {
			if service.UID == dst.ServiceUID {
				if service.Spec.ClusterIP == "" {
					log.Println("[Error]: Service has no cluster ip")
					return errors.New("service has no cluster ip")
				}

				serviceIP[index] = service.Spec.ClusterIP
				serviceExist = true
				break
			}
		}

		if !serviceExist {
			log.Println("[Error]: no corresponding service in dns")
			return nil
		}

		index++
	}

	log.Printf("[INFO]: Now, DNS %v(hostname is %v) is ready to start nginx docker\n", dns.UID, hostname)
	log.Printf("[INFO]: Paths number is %v", len(paths))
	containerID, err := pr.StartDNSNginxDocker(hostname, paths, serviceIP, port)
	if err != nil {
		log.Println("[Error]: DNS Nginx docker start error")
		return nil
	}

	log.Printf("[INFO]: DNS %v has been built, containerID is %v", dns.UID, containerID)
	pr.DNSMap[dns.UID] = DNSElement{ContainerID: containerID}
	return nil
}

func (pr *ProxyRuntime) DeleteDNS(dns *object.Dns) error {
	return nil
}

func (pr *ProxyRuntime) ModifyDNS(dns *object.Dns) error {
	return nil
}
