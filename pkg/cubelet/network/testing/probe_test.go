package testing

import (
	"Cubernetes/pkg/cubelet/network"
	"testing"
)

//const (
//	containerID   = "17b4889b4905"
//	netns         = "/var/run/netns/ns1"
//	podName       = "pod"
//	defaultIfName = "eth0"
//)

func TestDNSServer(t *testing.T) {
	network.InitDNSServer()
}
