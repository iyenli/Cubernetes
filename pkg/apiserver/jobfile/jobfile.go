package jobfile

import (
	cubeconfig "Cubernetes/config"
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func GetJobFile(JobUID string, filename string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/file/" + JobUID
	resp, err := http.Get(url)

	if err != nil {
		log.Println("fail to send http get request, err: ", err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP GET NOT OK when getting job file, status code: %v, server response: %s\n", resp.StatusCode, string(body))
		return errors.New(string(body))
	}

	err = ioutil.WriteFile(filename, body, 0777)
	if err != nil {
		log.Println("fail to write into file: ", filename, ", err: ", err)
		return err
	}

	return nil
}

func PostJobFile(JobUID string, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("fail to read file: ", filename, ", err: ", err)
		return err
	}

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Println("fail to create form file, err: ", err)
		return err
	}

	_, err = fileWriter.Write(buf)
	if err != nil {
		log.Println("fail to write form file, err: ", err)
		return err
	}

	err = writer.Close()
	if err != nil {
		log.Println("fail to close form file, err: ", err)
		return err
	}

	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/file/" + JobUID
	resp, err := http.Post(url, writer.FormDataContentType(), &body)

	if err != nil {
		log.Println("fail to send http post request, err: ", err)
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http post response body, err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP POST NOT OK when posting job file, status code: %v, server response: %s\n", resp.StatusCode, string(respBody))
		return errors.New(string(respBody))
	}

	return nil
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
