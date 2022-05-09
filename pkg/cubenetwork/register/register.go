package register

import cubeconfig "Cubernetes/config"

func RegisterToMaster(args []string) {
	if len(args) == 3 { // not same machine with api server
		cubeconfig.APIServerIp = args[2]
	}

	// TODO: Register
}

// Save uuid to log meta
func saveUUID(uuid string) error {
	return nil
}
