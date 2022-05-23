package objfile

import (
	cubeconfig "Cubernetes/config"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	resp, err := http.Get(url)

	if err != nil {
		log.Println("fail to send http get request, err: ", err)
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body, err: ", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP GET NOT OK when getting job output, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return "", errors.New(string(body))
	}

	return string(body), nil
}

func PostJobOutput(JobUID string, output string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/output/" + JobUID

	resp, err := http.Post(url, "text/plain", strings.NewReader(output))
	if err != nil || resp == nil {
		log.Println("fail to send http post request, err: ", err)
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http post response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP POST NOT OK when posting job output, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return errors.New(string(body))
	}

	return nil
}
