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

func GetActionFileStr(ScriptUID string) (string, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ScriptUID
	return getFileStr(url)
}

func PostActionFileStr(ScriptUID string, content string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/file/" + ScriptUID
	return postFileStr(url, content)
}
