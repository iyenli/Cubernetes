package objfile

import (
	cubeconfig "Cubernetes/config"
	"strconv"
)

func GetActionFile(ActionUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ActionUID
	return getFile(url, filename)
}

func PostActionFile(ActionUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ActionUID
	return postFile(url, filename)
}
