package proxyruntime

import (
	"Cubernetes/pkg/object"
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"log"
	"net"
	"strconv"
)

type PodTablesRuntime struct {
	ipt *iptables.IPTables
}

// InitPodChain Think twice and write pod's iptable
func InitPodChain() (*PodTablesRuntime, error) {
	ptr := &PodTablesRuntime{ipt: nil}
	err := ptr.InitObject()
	if err != nil {
		return nil, err
	}

	return ptr, nil
}

// InitObject private function! Just for test
func (ptr *PodTablesRuntime) InitObject() (err error) {
	ptr.ipt, err = iptables.New(iptables.Timeout(3))
	if err != nil {
		log.Println(err)
		return err
	}
	return
}

func ReleasePodChain() error {
	return nil
}

// AddPod FIX: -p should set front of --dport
func (ptr *PodTablesRuntime) AddPod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			err := ptr.ipt.Append(NatTable, PreRouting,
				"-d", pod.Status.IP.String(),
				"-p", port.Protocol,
				"--dport", strconv.FormatInt(int64(port.HostPort), 10),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", dockerIP.String(), strconv.FormatInt(int64(port.ContainerPort), 10)))

			if err != nil {
				log.Println("Add pod IP to iptables failed")
				return err
			}
		}
	}

	return nil
}

// DeletePod FIX: -p should set front of --dport
func (ptr *PodTablesRuntime) DeletePod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {

			err := ptr.ipt.DeleteIfExists(NatTable, PreRouting,
				"-d", pod.Status.IP.String(),
				"-p", port.Protocol,
				"--dport", strconv.FormatInt(int64(port.HostPort), 10),
				"-j", DnatOP,
				"--to-destination", fmt.Sprintf("%v:%v", dockerIP.String(), strconv.FormatInt(int64(port.ContainerPort), 10)))

			if err != nil {
				return err
			}
		}
	}

	return nil
}
