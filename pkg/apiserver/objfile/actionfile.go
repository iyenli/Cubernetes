package objfile

import (
	cubeconfig "Cubernetes/config"
	"strconv"
)

func GetActionFile(ActionName string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ActionName
	return getFile(url, filename)
}

func PostActionFile(ActionName string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ActionName
	return postFile(url, filename)
}
