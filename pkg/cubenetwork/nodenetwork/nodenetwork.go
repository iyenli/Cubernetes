package nodenetwork

import (
	cubeconfig "Cubernetes/config"
)

func SetMasterIP(IP string) {
	cubeconfig.APIServerIp = IP
}
