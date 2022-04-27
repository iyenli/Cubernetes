package testing

import (
	"Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/network"
	"testing"
)

func TestProbePlugin(t *testing.T) {
	var plugin network.CniNetworkPluginInterface
	plugin = network.ProbeNetworkPlugins("", "")
	host := NewFakeHost(make(map[string][]*network.PortMapping))
	err := plugin.Init(host, "", 0)

	if err != nil {
		return
	}
	id := container.ContainerID{
		Type: "default",
		ID:   "d329074c5ea4",
	}

	err = plugin.SetUpPod("default", "pause", id)
	if err != nil {
		return
	}

}
