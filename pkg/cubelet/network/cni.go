/*
	CNI Interface:
	type CNI interface {
		AddNetworkList(net *NetworkConfigList, rt *RuntimeConf) (types.Result, error)
		DelNetworkList(net *NetworkConfigList, rt *RuntimeConf) error

		AddNetwork(net *NetworkConfig, rt *RuntimeConf) (types.Result, error)
		DelNetwork(net *NetworkConfig, rt *RuntimeConf) error
	}

	This file wrap up cni interfaces.
*/
package network

import (
	"Cubernetes/pkg/cubelet/container"
	network "Cubernetes/pkg/cubelet/network/options"
	"context"
	"errors"
	"fmt"
	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
	"log"
	osexec "os/exec"
	"sort"
	"strings"
)

// ProbeNetworkPlugins use ["", ""] as paras would be replaced by default value
// They will be replaced with "/etc/cni/net.d" and "/opt/cni/bin"
func ProbeNetworkPlugins(pluginDir, binDir string) CniNetworkPluginInterface {
	if binDir == "" {
		binDir = network.DefaultCNIDir
	}

	plugin := &CniNetworkPlugin{
		// reserved for later use
		defaultNetwork: nil,
		// `bindir/loopback` must exist
		// Just support linux in Cubernetes:)
		// TODO: Check correctness
		loNetwork: getLoNetwork(binDir),
		pluginDir: pluginDir,
		binDir:    binDir,
	}

	// sync NetworkConfig in the best effort during probing.
	// probe and sync network config
	plugin.syncNetworkConfig()

	return plugin
}

// Probe network and set default network for plugin
func (plugin *CniNetworkPlugin) syncNetworkConfig() {
	defaultCNINetwork, err := getDefaultCNINetwork(plugin.pluginDir, plugin.binDir)
	if err != nil {
		log.Printf("Unable to update cni config: %s", err)
		return
	}
	plugin.setDefaultNetwork(defaultCNINetwork)
}

/* binDir == "/opt/cni/bin", ls binDir: (default cni plugin maintained)
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
* */
func getDefaultCNINetwork(pluginDir, binDir string) (*cniNetwork, error) {
	if pluginDir == "" {
		// net config directory
		pluginDir = network.DefaultNetDir
	}

	// search all configs under this directory, support .conf and .json
	files, err := libcni.ConfFiles(pluginDir, []string{".conf", ".json"})
	switch {
	case err != nil:
		return nil, err
	case len(files) == 0:
		return nil, fmt.Errorf("no networks found in %s", pluginDir)
	}

	sort.Strings(files)
	// Check config files until a legal one
	for _, confFile := range files {
		var confList *libcni.NetworkConfigList

		conf, err := libcni.ConfFromFile(confFile)
		// Check config file
		if err != nil {
			log.Printf("Error loading CNI config file %s: %v", confFile, err)
			continue
		}
		// Ensure the config has a "type" so we know what plugin to run.
		if conf.Network.Type == "" {
			log.Printf("Error loading CNI config file %s: no 'type';", confFile)
			continue
		}

		confList, err = libcni.ConfListFromConf(conf)
		if err != nil {
			log.Printf("Error converting CNI config file %s to list: %v", confFile, err)
			continue
		}

		if len(confList.Plugins) == 0 {
			log.Printf("CNI config list %s has no networks, skipping", confFile)
			continue
		}

		// TODO: Check correctness
		cniConfig := &libcni.CNIConfig{
			Path: []string{binDir},
		}
		res := &cniNetwork{name: confList.Name, NetworkConfig: confList, CNIConfig: cniConfig}
		return res, nil
	}

	return nil, fmt.Errorf("no valid networks found in %s", pluginDir)
}

func (plugin *CniNetworkPlugin) getDefaultNetwork() *cniNetwork {
	plugin.RLock()
	defer plugin.RUnlock()
	return plugin.defaultNetwork
}

func (plugin *CniNetworkPlugin) setDefaultNetwork(n *cniNetwork) {
	plugin.Lock()
	defer plugin.Unlock()
	plugin.defaultNetwork = n
}

func (plugin *CniNetworkPlugin) checkInitialized() error {
	if plugin.getDefaultNetwork() == nil {
		return errors.New("cni config uninitialized")
	}
	return nil
}

func (plugin *CniNetworkPlugin) Name() string {
	return network.CNIPluginName
}

func (plugin *CniNetworkPlugin) Status() error {
	// sync network config from pluginDir periodically to detect network config updates
	plugin.syncNetworkConfig()

	// Can't set up pods if we don't have any CNI network configs yet
	return plugin.checkInitialized()
}

// podSandboxID: Pod's pause container ID
func (plugin *CniNetworkPlugin) addToNetwork(network *cniNetwork, podName string, podNamespace string,
	podSandboxID container.ContainerID, podNetnsPath string) (types.Result, error) {
	rt, err := plugin.buildCNIRuntimeConf(podName, podNamespace, podSandboxID, podNetnsPath)
	if err != nil {
		log.Fatalf("Error adding network when building cni runtime conf: %v", err)
		return nil, err
	}

	netConf, cniNet := network.NetworkConfig, network.CNIConfig
	log.Printf("About to add CNI network %v (type=%v)", netConf.Name, netConf.Plugins[0].Network.Type)
	res, err := cniNet.AddNetworkList(context.TODO(), netConf, rt)
	if err != nil {
		log.Fatalf("Error adding network: %v", err)
		return nil, err
	}

	return res, nil
}

func (plugin *CniNetworkPlugin) deleteFromNetwork(network *cniNetwork, podName string, podNamespace string, podSandboxID container.ContainerID, podNetnsPath string) error {
	rt, err := plugin.buildCNIRuntimeConf(podName, podNamespace, podSandboxID, podNetnsPath)
	if err != nil {
		log.Fatalf("Error deleting network when building cni runtime conf: %v", err)
		return err
	}

	netConf, cniNet := network.NetworkConfig, network.CNIConfig
	log.Printf("About to del CNI network %v (type=%v)", netConf.Name, netConf.Plugins[0].Network.Type)
	err = cniNet.DelNetworkList(context.TODO(), netConf, rt)
	if err != nil {
		log.Fatalf("Error deleting network: %v", err)
		return err
	}
	return nil
}

func (plugin *CniNetworkPlugin) buildCNIRuntimeConf(podName string, podNs string, podSandboxID container.ContainerID, podNetnsPath string) (*libcni.RuntimeConf, error) {
	log.Printf("Got netns path %v", podNetnsPath)
	log.Printf("Using podns path %v", podNs)

	rt := &libcni.RuntimeConf{
		ContainerID: podSandboxID.ID,
		NetNS:       podNetnsPath,
		IfName:      network.DefaultInterfaceName,
		Args: [][2]string{
			{"IgnoreUnknown", "1"},
			{"POD_NAMESPACE", podNs},
			{"POD_NAME", podName},
			{"POD_INFRA_CONTAINER_ID", podSandboxID.ID},
		},
	}

	// port mappings are a cni capability-based args, rather than parameters
	// to a specific plugin
	portMappings, err := plugin.host.GetPodPortMappings(podSandboxID.ID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve port mappings: %v", err)
	}
	portMappingsParam := make([]CniPortMapping, 0, len(portMappings))
	for _, p := range portMappings {
		if p.HostPort <= 0 {
			continue
		}
		portMappingsParam = append(portMappingsParam, CniPortMapping{
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      strings.ToLower(p.Protocol),
			HostIP:        p.HostIP,
		})
	}
	rt.CapabilityArgs = map[string]interface{}{
		"portMappings": portMappingsParam,
	}

	return rt, nil
}

// Init Modify: cancel support for hairpin mode in bridge for simplicity
func (plugin *CniNetworkPlugin) Init(host Host, nonMasqueradeCIDR string, mtu int) error {
	var err error
	// in linux machine, install nsenter to maintain namespace
	plugin.nsEnterPath, err = osexec.LookPath("nsenter")
	if err != nil {
		return err
	}

	plugin.host = host
	plugin.syncNetworkConfig()
	return nil
}

func (plugin *CniNetworkPlugin) SetUpPod(namespace string, name string, id container.ContainerID) error {
	if err := plugin.checkInitialized(); err != nil {
		return err
	}
	netnsPath, err := plugin.host.GetNetNS(id.ID)
	if err != nil {
		return fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}

	if plugin.loNetwork != nil {
		if _, err = plugin.addToNetwork(plugin.loNetwork, name, namespace, id, netnsPath); err != nil {
			log.Fatalf("Error while adding to cni lo network: %s", err)
			return err
		}
	}

	_, err = plugin.addToNetwork(plugin.getDefaultNetwork(), name, namespace, id, netnsPath)
	if err != nil {
		log.Fatalf("Error while adding to cni network: %s", err)
		return err
	}

	return err
}

func (plugin *CniNetworkPlugin) TearDownPod(namespace string, name string, id container.ContainerID) error {
	if err := plugin.checkInitialized(); err != nil {
		return err
	}

	// Lack of namespace should not be fatal on teardown
	netnsPath, err := plugin.host.GetNetNS(id.ID)
	if err != nil {
		log.Printf("CNI failed to retrieve network namespace path: %v", err)
	}

	return plugin.deleteFromNetwork(plugin.getDefaultNetwork(), name, namespace, id, netnsPath)
}

func (plugin *CniNetworkPlugin) Capabilities() int32 {
	panic("unsupported by now:)")
}

func (plugin *CniNetworkPlugin) Event(name string, details map[string]interface{}) {
	panic("unsupported by now:)")
}
