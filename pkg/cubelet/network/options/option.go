package network

const (
	// DefaultCNIDir user set their plugin under it
	DefaultCNIDir          = "/opt/cni/bin"
	CNIPluginName          = "cni"
	DefaultNetDir          = "/etc/cni/net.d"
	DefaultCNIVersion      = "0.4.0"
	DefaultInterfaceName   = "eth0"
	DefaultInterfacePrefix = "eth"

	Resolve  = "resolvconf"
	BaseAddr = "/etc/resolvconf/resolv.conf.d/"
	BaseFile = BaseAddr + "base"
	HeadFile = BaseAddr + "head"

	HeadFileContent = `
search weave.local.root
nameserver 172.17.0.1
`
)
