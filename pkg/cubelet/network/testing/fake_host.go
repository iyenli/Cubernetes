package testing

import "Cubernetes/pkg/cubelet/network"

type fakeNetworkHost struct {
	FakePortMappingGetter
	fakeNamespaceGetter
}

type FakePortMappingGetter struct {
	PortMaps map[string][]*network.PortMapping
}

func (t *fakeNetworkHost) GetNetNS(containerID string) (string, error) {
	return "/var/run/netns/ns1", nil
}

func (t *fakeNetworkHost) GetPodPortMappings(containerID string) ([]*network.PortMapping, error) {
	return t.PortMaps[containerID], nil
}

type fakeNamespaceGetter struct {
	ns string
}

func NewFakeHost(ports map[string][]*network.PortMapping) *fakeNetworkHost {
	host := &fakeNetworkHost{
		FakePortMappingGetter{PortMaps: ports},
		fakeNamespaceGetter{ns: "/var/run/netns/ns1"},
	}

	return host
}
