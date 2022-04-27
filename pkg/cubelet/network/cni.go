/*
	CNI Interface:
	type CNI interface {
		AddNetworkList(net *NetworkConfigList, rt *RuntimeConf) (types.Result, error)
		DelNetworkList(net *NetworkConfigList, rt *RuntimeConf) error

		AddNetwork(net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
		DelNetwork(net *NetworkConfig, rt *RuntimeConf) error
	}

	binDir == "/opt/cni/bin", ls binDir: (default cni plugin maintained)
	├── bandwidth
	├── bridge
	├── dhcp
	├── firewall
	├── flannel
	├── host-device
	├── host-local
	├── ipvlan
	├── loopback
	├── macvlan
	├── portmap
	├── ptp
	├── sbr
	├── static
	├── tuning
	├── vlan
	└── vrf

	This file wrap up cni interfaces.
*/
package network

import (
	"Cubernetes/pkg/cubelet/container"
	network "Cubernetes/pkg/cubelet/network/options"
	"context"
	gocni "github.com/containerd/go-cni"
	"log"
)

// ProbeNetworkPlugins use ["", ""] as paras would be replaced by default value
// They will be replaced with "/etc/cni/net.d" and "/opt/cni/bin"
func ProbeNetworkPlugins(pluginDir, binDir string) gocni.CNI {
	if binDir == "" {
		binDir = network.DefaultCNIDir
	}
	if pluginDir == "" {
		pluginDir = network.DefaultNetDir
	}

	plugin, err := gocni.New(
		gocni.WithMinNetworkCount(2),
		gocni.WithPluginConfDir("/etc/cni/net.d"),
		gocni.WithPluginDir([]string{"/opt/cni/bin"}),
		// Sets the prefix for network interfaces, eth by default
		gocni.WithInterfacePrefix(network.DefaultInterfacePrefix))
	if err != nil {
		log.Fatalf("failed to initialize cni library: %v", err)
	}

	// Load the cni configuration
	// Get lo and default network(warp up origin code)
	if err := plugin.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
		log.Fatalf("failed to load cni configuration: %v", err)
	}

	return plugin
}

func SetUpPod(cni gocni.CNI, netNamespace string, name string, id container.ContainerID) (*gocni.Result, error) {
	ctx := context.Background()

	// Setup network for namespace.
	labels := map[string]string{
		"K8S_POD_NAME":               name,
		"K8S_POD_INFRA_CONTAINER_ID": id.ID,
		// Plugin tolerates all Args embedded by unknown labels, like
		// K8S_POD_NAMESPACE/NAME/INFRA_CONTAINER_ID...
		"IgnoreUnknown": "1",
	}

	result, err := cni.Setup(ctx, id.ID, netNamespace, gocni.WithLabels(labels))
	if err != nil {
		log.Printf("failed to setup network for namespace: %v", err)
		return nil, err
	}

	return result, nil
}

func TearDownPod(cni gocni.CNI, netNamespace string, name string, id container.ContainerID) error {
	ctx := context.Background()

	labels := map[string]string{
		"K8S_POD_NAME":               name,
		"K8S_POD_INFRA_CONTAINER_ID": id.ID,
		"IgnoreUnknown":              "1",
	}

	if err := cni.Remove(ctx, id.ID, netNamespace, gocni.WithLabels(labels)); err != nil {
		log.Printf("failed to teardown network: %v", err)
		return err
	}

	return nil
}

func CheckPodStatus(cni gocni.CNI, netNamespace string, id container.ContainerID) error {
	ctx := context.Background()

	err := cni.Check(ctx, id.ID, netNamespace)
	if err != nil {
		return err
	}
	return nil
}
