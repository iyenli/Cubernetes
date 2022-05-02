package testing

import (
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

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
