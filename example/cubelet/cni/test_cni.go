package main

import (
	"context"
	"fmt"
	gocni "github.com/containerd/go-cni"
	"log"
)

func main() {
	id := "d329074c5ea4"
	netns := "/var/run/netns/ns1"

	// CNI allows multiple CNI configurations and the network interface
	// will be named by eth0, eth1, ..., ethN.
	ifPrefixName := "eth"
	defaultIfName := "eth0"

	// Initializes library
	l, err := gocni.New(
		// one for loopback network interface
		gocni.WithMinNetworkCount(2),
		gocni.WithPluginConfDir("/etc/cni/net.d"),
		gocni.WithPluginDir([]string{"/opt/cni/bin"}),
		// Sets the prefix for network interfaces, eth by default
		gocni.WithInterfacePrefix(ifPrefixName))
	if err != nil {
		log.Fatalf("failed to initialize cni library: %v", err)
	}

	// Load the cni configuration
	if err := l.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
		log.Fatalf("failed to load cni configuration: %v", err)
	}

	// Setup network for namespace.
	labels := map[string]string{
		"K8S_POD_NAMESPACE":          "namespace1",
		"K8S_POD_NAME":               "pod1",
		"K8S_POD_INFRA_CONTAINER_ID": id,
		// Plugin tolerates all Args embedded by unknown labels, like
		// K8S_POD_NAMESPACE/NAME/INFRA_CONTAINER_ID...
		"IgnoreUnknown": "1",
	}

	ctx := context.Background()

	// Teardown network
	defer func() {
		if err := l.Remove(ctx, id, netns, gocni.WithLabels(labels)); err != nil {
			log.Fatalf("failed to teardown network: %v", err)
		}
	}()

	// Setup network
	result, err := l.Setup(ctx, id, netns, gocni.WithLabels(labels))
	if err != nil {
		log.Fatalf("failed to setup network for namespace: %v", err)
	}

	// Get IP of the default interface
	IP := result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
	fmt.Printf("IP of the default interface %s:%s", defaultIfName, IP)
}
