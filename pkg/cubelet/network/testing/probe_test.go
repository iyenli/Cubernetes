package testing

import (
	"Cubernetes/pkg/cubelet/container"
	"Cubernetes/pkg/cubelet/network"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	containerID   = "17b4889b4905"
	netns         = "/var/run/netns/ns1"
	podName       = "pod"
	defaultIfName = "eth0"
)

func TestCni(t *testing.T) {
	cni := network.ProbeNetworkPlugins("", "")

	result, err := network.SetUpPod(cni, netns, podName, container.ContainerID{
		Type: "common",
		ID:   containerID,
	})

	assert.NoError(t, err, "Setup without error")
	assert.Equal(t, len(result.Routes), 1)

	err = network.TearDownPod(cni, netns, podName, container.ContainerID{
		Type: "common",
		ID:   containerID,
	})
	assert.NoError(t, err, "Teardown without error")

	// Get IP of the default interface
	fmt.Printf("CNI Create result: %v\n", result)
	IP := result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
	fmt.Printf("IP of the default interface %s:%s \n", defaultIfName, IP)
}
