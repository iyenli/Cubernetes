package health

import (
	cubeconfig "Cubernetes/config"
	"net/http"
	"strconv"
)

func CheckApiServerHealth() bool {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/health"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
