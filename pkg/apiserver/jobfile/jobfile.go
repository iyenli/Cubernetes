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
		log.Println("fail to send http get request")
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("job file not found")
	}

	err = ioutil.WriteFile(filename, body, 0777)
	if err != nil {
		log.Println("fail to write into file: ", filename)
		return err
	}

	return nil
}

func PostJobFile(JobUID string, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("fail to read file: ", filename)
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

	if err != nil || resp == nil {
		log.Println("fail to send http post request")
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New("fail to upload job file")
	}

	return nil
}

func GetJobOutput(JobUID string) (string, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/output/" + JobUID
	resp, err := http.Get(url)

	if err != nil {
		log.Println("fail to send http get request")
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read http get response body")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("job not finished")
	}

	return string(body), nil
}

func PostJobOutput(JobUID string, output string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/output/" + JobUID

	resp, err := http.Post(url, "text/plain", strings.NewReader(output))
	if err != nil || resp == nil {
		log.Println("fail to send http post request")
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New("fail to post job output")
	}

	return nil
}
