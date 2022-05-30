package objfile

import (
	cubeconfig "Cubernetes/config"
	"strconv"
)

func GetActionFile(ScriptUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ScriptUID
	return getFile(url, filename)
}

func PostActionFile(ScriptUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ScriptUID
	return postFile(url, filename)
}
