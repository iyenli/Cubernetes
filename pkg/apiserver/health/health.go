package health

import (
	cubeconfig "Cubernetes/config"
	"log"
	"net/http"
	"strconv"
)

func CheckApiServerHealth() bool {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/health"
	resp, err := http.Get(url)

	if err != nil {
		log.Println("fail to send http get request")
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
