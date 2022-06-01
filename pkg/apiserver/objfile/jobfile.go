package objfile

import (
	cubeconfig "Cubernetes/config"
	"strconv"
)

func GetJobFile(JobUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/file/" + JobUID
	return getFile(url, filename)
}

func PostJobFile(JobUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/file/" + JobUID
	return postFile(url, filename)
}

func GetJobOutput(JobUID string) (string, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/output/" + JobUID
	return getFileStr(url)
}

func PostJobOutput(JobUID string, output string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/output/" + JobUID
	return postFileStr(url, output)
}
