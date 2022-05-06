package testing

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins/host"
	"Cubernetes/pkg/cubenetwork/weaveplugins/weave"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

// TestAddNode JCloud env
func TestAddNode(t *testing.T) {
	host1 := host.Host{IP: net.ParseIP("192.168.1.9")}
	host2 := host.Host{IP: net.ParseIP("192.168.1.5")}

	err := host.InitWeave(host1)
	assert.NoError(t, err)

	err = host.AddNode(host1, host2)
	assert.NoError(t, err)
}

// TestInitNode WSL env
func TestInitNode(t *testing.T) {
	host1 := host.Host{IP: net.ParseIP("172.17.0.19")}

	err := host.InitWeave(host1)
	assert.NoError(t, err)
}

func TestInstallWeave(t *testing.T) {
	err := weave.InstallWeave()
	assert.NoError(t, err)
}
