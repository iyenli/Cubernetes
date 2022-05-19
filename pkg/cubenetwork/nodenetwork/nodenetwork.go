package nodenetwork

import (
	cubeconfig "Cubernetes/config"
	"log"
)

func SetMasterIP(IP string) {
	cubeconfig.APIServerIp = IP
	log.Println("[INFO]: Set master IP:", IP)
}
