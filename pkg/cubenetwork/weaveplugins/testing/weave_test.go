package testing

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

// TestAddNode JCloud env
func TestAddNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	host1 := weaveplugins.Host{IP: net.ParseIP("192.168.1.9")}
	host2 := weaveplugins.Host{IP: net.ParseIP("192.168.1.5")}

	err = weaveplugins.AddNode(host1, host2)
	assert.NoError(t, err)
}

// TestInitNode WSL env
func TestInitNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	err = weaveplugins.InitWeave()
	assert.NoError(t, err)
}

func TestInstallWeave(t *testing.T) {
	err := weaveplugins.InstallWeave()
	assert.NoError(t, err)
}

func TestWeaveStatus(t *testing.T) {
	output, err := weaveplugins.CheckPeers()
	assert.NoError(t, err)

	t.Log(string(output))
}

func TestWeaveStop(t *testing.T) {
	err := weaveplugins.CloseNetwork()
	assert.NoError(t, err)
}

func TestAddPod(t *testing.T) {
	id := RunContainer()
	network, err := weaveplugins.AddPodToNetwork(id)
	assert.NoError(t, err)

	t.Logf("Container ID: %v, IP: %v", id, network.String())
	err = weaveplugins.DeletePodFromNetwork(id)
	assert.NoError(t, err)
}

func TestParseIP(t *testing.T) {
	ip := net.ParseIP("10.40.0.0")
	assert.NotNil(t, ip)
}

func TestDNSEntry(t *testing.T) {
	id := RunContainer()
	newString := uuid.NewString()[:5]

	err := weaveplugins.AddDNSEntry(newString, id)
	assert.NoError(t, err)

	err = weaveplugins.DeleteDNSEntry(id)
	assert.NoError(t, err)
}
