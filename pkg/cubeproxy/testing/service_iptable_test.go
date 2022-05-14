package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"Cubernetes/pkg/object"
	"log"
	"testing"
	"time"

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
				TargetPort: 8098,
			},
			{
				Protocol:   "UDP",
				Port:       80,
				TargetPort: 9841,
			},
		},
		ClusterIP: "10.0.1.3",
	},
	Status: &object.ServiceStatus{
		Ingress: nil,
	},
}

/**
iptables pkg test
*/

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

func TestRelease(t *testing.T) {
	rt, err := proxyruntime.InitProxyRuntime()
	assert.NoError(t, err)

	err = rt.ReleaseIPTables()
	assert.NoError(t, err)
}

// If test failed for
func TestClearIPT(t *testing.T) {
	ipt, err := iptables.New(iptables.Timeout(3))
	assert.NoError(t, err)

	err = ipt.ClearAll()
	assert.NoError(t, err)
	err = ipt.DeleteAll()
	assert.NoError(t, err)
}

/**
iptables pkg test end
*/

func TestInitIPTables(t *testing.T) {
	rt, err := proxyruntime.InitProxyRuntime()
	assert.NoError(t, err)

	// Check IP Tables and test release!
	//time.Sleep(10 * time.Second)

	err = rt.ReleaseIPTables()
	assert.NoError(t, err)
}

func TestAddService(t *testing.T) {
	rt, err := proxyruntime.InitProxyRuntime()
	assert.NoError(t, err)
	assert.NotNil(t, rt)

	err = rt.AddService(&service)
	assert.NoError(t, err)

	time.Sleep(40 * time.Second)
	// clear testing env
	defer func(rt *proxyruntime.ProxyRuntime) {
		err := rt.DeleteService(&service)
		if err != nil {
			t.Log(err)
			log.Panicln("Release proxy runtime failed.")
		}

		err = rt.ReleaseIPTables()
		if err != nil {
			t.Log(err)
			log.Panicln("Release proxy runtime failed.")
		}
	}(rt)
}

func TestDeleteService(t *testing.T) {
	rt, err := proxyruntime.InitProxyRuntime()
	assert.NoError(t, err)

	err = rt.DeleteService(&service)
	assert.NoError(t, err)

	// clear testing env
	defer func(rt *proxyruntime.ProxyRuntime) {
		err := rt.ReleaseIPTables()
		if err != nil {
			log.Panicln("Release proxy runtime failed.")
		}
	}(rt)
}
