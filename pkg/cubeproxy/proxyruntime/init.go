package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubeproxy/informer"
	"Cubernetes/pkg/object"
	"log"
)

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

	// Delete nginx config folder
	// For now, try what if not restart all nginx docker?
	//if _, err := os.Stat(options.NginxFile); err == nil {
	//	err := os.RemoveAll(options.NginxFile)
	//	if err != nil {
	//		log.Println("[Error]: clear nginx file failed")
	//		return pr, nil
	//	}
	//}

	return pr, nil
}

// AddAllExistService construct exist service iptables
func (pr *ProxyRuntime) AddAllExistService() error {
	services, err := crudobj.GetServices()
	if err != nil {
		log.Fatal("[Fatal]: Get services failed")
		return err
	}

	if len(pr.ServiceChainMap) != 0 {
		log.Println("[BUG]: Add exist service should be called when initializing")
	}
	err = pr.ServiceInformer.InitInformer(services)
	if err != nil {
		log.Fatalln("[Fatal]: Init pod informer failed")
	}

	for _, service := range services {
		log.Println("[INFO]: Cubeproxy init, add exist service, UID:", service.UID)
		err := pr.AddExistService(&service)
		if err != nil {
			log.Printf("[INFO]: Add exist service %v failed", service.UID)
		}
	}

	return nil
}

// AddExistService Service has correct endpoints, just add iptables mapping:)
func (pr *ProxyRuntime) AddExistService(service *object.Service) error {
	if service.Status == nil || len(service.Status.Endpoints) == 0 {
		log.Println("[INFO]: Adding a service without endpoint, do nothing")
		return nil
	}

	prob := make([][]string, len(service.Spec.Ports))
	for idx, _ := range prob {
		prob[idx] = make([]string, len(service.Status.Endpoints))
	}

	pr.ServiceChainMap[service.UID] = ServiceChainElement{
		serviceChainUid:     make([]string, len(service.Spec.Ports)),
		probabilityChainUid: prob,
		numberOfPods:        len(service.Status.Endpoints),
	}

	podIPs := make([]string, len(service.Status.Endpoints))
	for idx, end := range service.Status.Endpoints {
		podIPs[idx] = end.String()
	}

	for idx, port := range service.Spec.Ports {
		err := pr.MapPortToPods(service, podIPs, &port, idx)
		if err != nil {
			log.Println("[error]: map port to pods failed")
			return err
		}
	}
	return nil
}
