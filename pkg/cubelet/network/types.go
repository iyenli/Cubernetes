package network

import (
	"github.com/containernetworking/cni/libcni"
	"net"
	"sync"
)

// Host A host has a namespace and map of port
type Host interface {
	// NamespaceGetter is a getter for sandbox namespace information.
	NamespaceGetter

	// PortMappingGetter is a getter for sandbox port mapping information.
	PortMappingGetter
}

// PodPortMapping represents a pod's network state
type PodPortMapping struct {
	// pod's Namespace and mapping name
	Namespace string
	Name      string
	// pod's IP
	IP           net.IP
	PortMappings []*PortMapping
}

// PortMapping represents a network port in a container
// No need for container name: containers in same pod are in lo
type PortMapping struct {
	// pod name
	Name          string
	HostPort      int32
	ContainerPort int32
	Protocol      string
	HostIP        string
}

type NamespaceGetter interface {
	// GetNetNS returns network namespace information for the given containerID.
	// Runtimes should *never* return an empty namespace and nil error for
	// a container; if error is nil then the namespace string must be valid.
	GetNetNS(containerID string) (string, error)
}

type PortMappingGetter interface {
	// GetPodPortMappings returns sandbox port mappings information.
	GetPodPortMappings(containerID string) ([]*PortMapping, error)
}

// CniNetworkPlugin support 2 network, one of them is lo
type CniNetworkPlugin struct {
	loNetwork      *cniNetwork // local loop
	defaultNetwork *cniNetwork

	host        Host
	nsEnterPath string
	pluginDir   string
	binDir      string

	sync.RWMutex
}

// cniNetwork private struct
type cniNetwork struct {
	name          string
	NetworkConfig *libcni.NetworkConfigList
	CNIConfig     libcni.CNI
}

// CniPortMapping maps to CNI port mapping
type CniPortMapping struct {
	HostPort      int32  `json:"hostPort"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP"`
}

// CniNetworkPluginInterface NetworkPlugin Plugin is an interface to network plugins for the kubelet
type CniNetworkPluginInterface interface {
	// Init initializes the plugin. This will be called exactly once
	// before any other methods are called.
	Init(host Host, nonMasqueradeCIDR string, mtu int) error

	Event(name string, details map[string]interface{})

	// Name returns the plugin's name.
	Name() string

	// Capabilities Returns a set of NET_PLUGIN_CAPABILITY_*
	Capabilities() int32

	// SetUpPod is the method called after the infra container of
	// the pod has been created but before the other containers of the
	// pod are launched.
	SetUpPod(namespace string, name string, podSandboxID string, annotations map[string]string) error

	// TearDownPod is the method called before a pod's infra container will be deleted
	TearDownPod(namespace string, name string, podSandboxID int) error

	// GetPodNetworkStatus is the method called to obtain the ipv4 or ipv6 addresses of the container
	GetPodNetworkStatus(namespace string, name string, podSandboxID int) (int, error)

	// Status returns error if the network plugin is in error state
	Status() error
}
