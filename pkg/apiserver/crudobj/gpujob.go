package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetGpuJob(UID string) (object.GpuJob, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.GpuJob{}, err
	}

	var job object.GpuJob
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Println("fail to parse GpuJob")
		return object.GpuJob{}, err
	}

	return job, nil
}

func GetGpuJobs() ([]object.GpuJob, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJobs"

	body, err := getRequest(url)
	if err != nil {
		log.Println("[Error]: getRequest fail")
		return nil, err
	}

	var jobs []object.GpuJob
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		log.Println("[Warn]: fail to parse GpuJobs, output:", string(body))
		return nil, err
	}

	return jobs, nil
}

func SelectGpuJobs(selectors map[string]string) ([]object.GpuJob, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/gpuJobs"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var jobs []object.GpuJob
	err = json.Unmarshal(body, &jobs)
	if err != nil {
		log.Println("fail to parse GpuJobs")
		return nil, err
	}

	return jobs, nil
}

func CreateGpuJob(job object.GpuJob) (object.GpuJob, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob"

	body, err := postRequest(url, job)
	if err != nil {
		log.Println("postRequest fail")
		return job, err
	}

	var newJob object.GpuJob
	err = json.Unmarshal(body, &newJob)
	if err != nil {
		log.Println("fail to parse GpuJob")
		return job, err
	}

	return newJob, nil
}

func UpdateGpuJob(job object.GpuJob) (object.GpuJob, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/" + job.UID

	body, err := putRequest(url, job)
	if err != nil {
		log.Println("putRequest fail")
		return job, err
	}

	var newJob object.GpuJob
	err = json.Unmarshal(body, &newJob)
	if err != nil {
		log.Println("fail to parse GpuJob")
		return job, err
	}

	return newJob, nil
}

func DeleteGpuJob(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/gpuJob/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
