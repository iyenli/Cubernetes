package options

const (
	containerdRuntimeEndpoint = "/run/containerd/containerd.sock"
)

type CubeletFlags struct {
	// The Cubelet will load its initial configuration from this file.
	// Omit this flag to use the combination of built-in default configuration values and flags.
	CubeletConfigFile string
	// remoteRuntimeEndpoint is the endpoint of remote runtime service
	RemoteRuntimeEndpoint string
	// remoteImageEndpoint is the endpoint of remote image service
	RemoteImageEndpoint string
}

func NewCubeletFlags() *CubeletFlags {

	return &CubeletFlags{
		RemoteRuntimeEndpoint: containerdRuntimeEndpoint,
		RemoteImageEndpoint:   containerdRuntimeEndpoint,
	}
}
