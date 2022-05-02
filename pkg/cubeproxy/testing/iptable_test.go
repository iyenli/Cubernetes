package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"Cubernetes/pkg/object"
	"net"
	"testing"

	"github.com/coreos/go-iptables/iptables"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	var flag bool
	var err error

	ipTable, err := iptables.New(iptables.Timeout(3))
	assert.NoError(t, err)

	flag, err = ipTable.ChainExists(proxyruntime.FilterTable, proxyruntime.InputChain)
	assert.True(t, flag)
	assert.NoError(t, err)

	flag, err = ipTable.ChainExists(proxyruntime.NatTable, proxyruntime.OutputChain)
	assert.True(t, flag)
	assert.NoError(t, err)

	flag, err = ipTable.ChainExists(proxyruntime.NatTable, proxyruntime.DockerChain)
	assert.True(t, flag)
	assert.NoError(t, err)

	flag, err = ipTable.ChainExists(proxyruntime.FilterTable, proxyruntime.DockerChain)
	assert.True(t, flag)
	assert.NoError(t, err)

	flag, err = ipTable.ChainExists(proxyruntime.NatTable, "Faker")
	assert.False(t, flag)
	assert.NoError(t, err)
}

func TestChainDelete(t *testing.T) {
	ipTable, err := iptables.New(iptables.Timeout(3))
	assert.NoError(t, err)

	err = ipTable.DeleteIfExists(proxyruntime.NatTable, "NewChain")
	assert.NoError(t, err)
}

func TestAddService(t *testing.T) {
	err := proxyruntime.InitIPTables()
	assert.NoError(t, err)

	service := object.Service{
		TypeMeta: object.TypeMeta{
			Kind:       "Service",
			APIVersion: "C8s/v1",
		},
		ObjectMeta: object.ObjectMeta{
			Name:      "my-service",
			Namespace: "ns",
			UID:       "pp-qq",
		},
		Spec: object.ServiceSpec{
			Selector: map[string]string{
				"gpu": "Nvidia",
			},
			Ports: []object.ServicePort{
				{
					Protocol:   "TCP",
					Port:       8080,
					TargetPort: 8080,
				},
				{
					Protocol:   "UDP",
					Port:       9870,
					TargetPort: 9841,
				},
			},
			ClusterIP: "10.0.1.3",
		},
		Status: &object.ServiceStatus{},
	}

	err = proxyruntime.AddService(&service)
	assert.NoError(t, err)
}

func TestAddPod(t *testing.T) {
	err := proxyruntime.InitIPTables()
	assert.NoError(t, err)

	pod, err := proxyruntime.GetPodByService(nil)
	assert.NoError(t, err)
	assert.NotNil(t, pod)
	assert.Equal(t, 1, len(pod))

	containerIP := net.ParseIP("10.0.0.4")
	t.Log(pod[0].Status.IP.String())
	err = proxyruntime.AddPod(&pod[0], containerIP)
	assert.NoError(t, err)

	// Now check IP Table manually...
}

func TestDeletePod(t *testing.T) {
	err := proxyruntime.InitIPTables()
	assert.NoError(t, err)

	pod, err := proxyruntime.GetPodByService(nil)
	assert.NoError(t, err)
	assert.NotNil(t, pod)
	assert.Equal(t, 1, len(pod))

	containerIP := net.ParseIP("10.0.0.4")
	err = proxyruntime.DeletePod(&pod[0], containerIP)
	assert.NoError(t, err)

	// Now check IP Table manually...
}
