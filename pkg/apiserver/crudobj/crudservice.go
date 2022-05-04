package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetService(UID string) (object.Service, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/service/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Service{}, err
	}

	var service object.Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		log.Println("fail to parse Service")
		return object.Service{}, err
	}

	return service, nil
}

func GetServices() ([]object.Service, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/services"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var services []object.Service
	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Println("fail to parse Services")
		return nil, err
	}

	return services, nil
}

func SelectServices(selectors map[string]string) ([]object.Service, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/services"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var services []object.Service
	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Println("fail to parse Services")
		return nil, err
	}

	return services, nil
}

func CreateService(service object.Service) (object.Service, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/service"

	body, err := postRequest(url, service)
	if err != nil {
		log.Println("postRequest fail")
		return service, err
	}

	var newService object.Service
	err = json.Unmarshal(body, &newService)
	if err != nil {
		log.Println("fail to parse Service")
		return service, err
	}

	return newService, nil
}

func UpdateService(service object.Pod) (object.Pod, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/service/" + service.UID

	body, err := putRequest(url, service)
	if err != nil {
		log.Println("putRequest fail")
		return service, err
	}

	var newService object.Pod
	err = json.Unmarshal(body, &newService)
	if err != nil {
		log.Println("fail to parse Service")
		return service, err
	}

	return newService, nil
}

func DeleteService(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/service/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
