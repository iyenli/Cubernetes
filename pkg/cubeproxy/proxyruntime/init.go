package proxyruntime

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubelet/dockershim"
	"Cubernetes/pkg/cubeproxy/informer"
	"Cubernetes/pkg/object"
	"github.com/coreos/go-iptables/iptables"
	"log"
	"net"
)

func InitProxyRuntime() (*ProxyRuntime, error) {
	dockerInstance, err := dockershim.NewDockerRuntime()
	if err != nil {
		log.Println("[Error]: Init docker runtime error")
	}

	pr := &ProxyRuntime{
		Ipt:            nil,
		DockerInstance: dockerInstance,

		ServiceInformer: informer.NewServiceInformer(),
		PodInformer:     informer.NewPodInformer(),
		DNSInformer:     informer.NewDNSInformer(),

		ServiceChainMap: make(map[string]ServiceChainElement),
		DNSMap:          make(map[string]DNSElement),
	}

	err = pr.InitObject()
	if err != nil {
		log.Panicln("[Warn]: Init object failed")
		return nil, err
	}

	/* check env */
	flag, err := pr.Ipt.ChainExists(FilterTable, DockerChain)
	if !flag {
		log.Printf("[Warn]: Start docker first")
		//return nil, err
	}
	flag, err = pr.Ipt.ChainExists(NatTable, DockerChain)
	if !flag {
		log.Printf("[Warn]: Start docker first")
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
	if service.Spec.ClusterIP == "" || len(service.Spec.Selector) == 0 {
		log.Println("[INFO]: Service without cluster ip or selector, ignore")
		return nil
	}

	if service.Status == nil {
		service.Status = &object.ServiceStatus{
			Endpoints: []net.IP{},
			Ingress:   []object.PodIngress{},
		}
	}

	service.Status.Endpoints = []net.IP{}
	pods := pr.PodInformer.ListPods()
	for _, pod := range pods {
		if object.MatchLabelSelector(service.Spec.Selector, pod.Labels) {
			if pod.Status != nil && pod.Status.IP != nil {
				service.Status.Endpoints = append(service.Status.Endpoints, pod.Status.IP)
			}
		}
	}

	if service.Status == nil || len(service.Status.Endpoints) == 0 {
		log.Println("[INFO]: Adding a service without endpoint, do nothing")
		return nil
	}

	prob := make([][]string, len(service.Spec.Ports))
	for idx, _ := range prob {
		prob[idx] = make([]string, len(service.Status.Endpoints))
	}

	pr.ServiceChainMap[service.UID] = ServiceChainElement{
		ServiceChainUid:     make([]string, len(service.Spec.Ports)),
		ProbabilityChainUid: prob,
		NumberOfPods:        0,
	}

	podIPs := make([]string, len(service.Status.Endpoints))
	for idx, end := range service.Status.Endpoints {
		podIPs[idx] = end.String()
	}

	for idx, port := range service.Spec.Ports {
		err := pr.MapPortToPods(service, podIPs, &port, idx)
		if err != nil {
			log.Println("[error]: map0 port to pods failed")
			return err
		}
	}

	_, err := crudobj.UpdateService(*service)
	if err != nil {
		log.Println("[Error]: update service failed")
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

func (pr *ProxyRuntime) ClearAllService() error {
	for _, service := range pr.ServiceInformer.ListServices() {
		err := pr.DeleteService(&service)
		if err != nil {
			return err
		}
	}

	return nil
}
