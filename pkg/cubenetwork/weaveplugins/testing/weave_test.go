package testing

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/dns"
	"Cubernetes/pkg/cubenetwork/weaveplugins/host"
	"Cubernetes/pkg/cubenetwork/weaveplugins/pod"
	"Cubernetes/pkg/cubenetwork/weaveplugins/weave"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

// TestAddNode JCloud env
func TestAddNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	host1 := host.Host{IP: net.ParseIP("192.168.1.9")}
	host2 := host.Host{IP: net.ParseIP("192.168.1.5")}

	err = host.AddNode(host1, host2)
	assert.NoError(t, err)
}

// TestInitNode WSL env
func TestInitNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	err = host.InitWeave()
	assert.NoError(t, err)
}

func TestInstallWeave(t *testing.T) {
	err := weave.InstallWeave()
	assert.NoError(t, err)
}

func TestWeaveStatus(t *testing.T) {
	output, err := host.CheckPeers()
	assert.NoError(t, err)

	t.Log(string(output))
}

func TestWeaveStop(t *testing.T) {
	err := host.CloseNetwork()
	assert.NoError(t, err)
}

func TestAddPod(t *testing.T) {
	id := RunContainer()
	network, err := pod.AddPodToNetwork(id)
	assert.NoError(t, err)

	t.Logf("Container ID: %v, IP: %v", id, network.String())
	err = pod.DeletePodFromNetwork(id)
	assert.NoError(t, err)
}

func TestParseIP(t *testing.T) {
	ip := net.ParseIP("10.40.0.0")
	assert.NotNil(t, ip)
}

func TestDNSEntry(t *testing.T) {
	id := RunContainer()
	newString := uuid.NewString()[:5]

	err := dns.AddDNSEntry(newString, id)
	assert.NoError(t, err)

	err = dns.DeleteDNSEntry(id)
	assert.NoError(t, err)
}
