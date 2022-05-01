package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"github.com/coreos/go-iptables/iptables"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestIPAppend(t *testing.T) {
	ipTable, err := iptables.New(iptables.Timeout(3))
	assert.NoError(t, err)

	err = ipTable.DeleteChain(proxyruntime.NatTable, "NewChain")
	assert.NoError(t, err)
}
