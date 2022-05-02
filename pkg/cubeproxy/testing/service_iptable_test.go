package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"Cubernetes/pkg/object"
	"testing"

	"github.com/coreos/go-iptables/iptables"
	"github.com/stretchr/testify/assert"
)

var service = object.Service{
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

func TestInitIPTables(t *testing.T) {
	err := proxyruntime.InitIPTables()
	assert.NoError(t, err)
}

func TestReleaseIPTables(t *testing.T) {
	err := proxyruntime.InitObject()
	err = proxyruntime.ReleaseIPTables()
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

	err = proxyruntime.AddService(&service)
	assert.NoError(t, err)
}

func TestDeleteService(t *testing.T) {
	err := proxyruntime.InitIPTables()
	assert.NoError(t, err)

	err = proxyruntime.DeleteService(&service)
	assert.NoError(t, err)
}
