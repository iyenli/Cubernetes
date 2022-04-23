package network

import "github.com/containernetworking/cni/libcni"

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
