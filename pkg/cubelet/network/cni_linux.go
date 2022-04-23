package network

import (
	"Cubernetes/pkg/cubelet/container"
	network "Cubernetes/pkg/cubelet/network/options"
	"fmt"
	"github.com/containernetworking/cni/libcni"
)

func getLoNetwork(binDir string) *cniNetwork {
	loConfig, err := libcni.ConfListFromBytes([]byte(`{
  		"cniVersion": "0.4.0",
  		"name": "cni-loopback",
  		"plugins":[{
    		"type": "loopback"
  		}]
	}`))

	if err != nil {
		panic(err)
	}
	cniConfig := &libcni.CNIConfig{
		Path: []string{binDir},
	}

	// first network: local loop
	loNetwork := &cniNetwork{
		name:          "lo",
		NetworkConfig: loConfig,
		CNIConfig:     cniConfig,
	}
	return loNetwork
}

func (plugin *CniNetworkPlugin) GetPodNetworkStatus(namespace string, name string, id container.ContainerID) (*PodNetworkStatus, error) {
	netnsPath, err := plugin.host.GetNetNS(id)
	if err != nil {
		return nil, fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}
	if netnsPath == "" {
		return nil, fmt.Errorf("cannot find the network namespace, skipping pod network status for container %q", id)
	}

	ip, err := GetPodIP(plugin.nsEnterPath, netnsPath, network.DefaultInterfaceName)
	if err != nil {
		return nil, err
	}

	return &PodNetworkStatus{IP: ip}, nil
}
