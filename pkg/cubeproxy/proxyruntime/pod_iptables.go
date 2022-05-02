package proxyruntime

import (
	"Cubernetes/pkg/object"
	"fmt"
	"log"
	"net"
	"strconv"
)

// InitPodChain Think twice and write pod's iptable
func InitPodChain() error {
	return nil
}

func ReleasePodChain() error {
	return nil
}

// AddPod FIX: -p should set front of --dport
func AddPod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			err := ipt.Append(NatTable, PreRouting,
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
func DeletePod(pod *object.Pod, dockerIP net.IP) error {
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {

			err := ipt.DeleteIfExists(NatTable, PreRouting,
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
