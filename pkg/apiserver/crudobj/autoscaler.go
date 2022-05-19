package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetAutoScaler(UID string) (object.AutoScaler, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/autoScaler/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.AutoScaler{}, err
	}

	var as object.AutoScaler
	err = json.Unmarshal(body, &as)
	if err != nil {
		log.Println("fail to parse AutoScaler")
		return object.AutoScaler{}, err
	}

	return as, nil
}

func GetAutoScalers() ([]object.AutoScaler, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/autoScalers"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var autoScalers []object.AutoScaler
	err = json.Unmarshal(body, &autoScalers)
	if err != nil {
		log.Println("fail to parse AutoScalers")
		return nil, err
	}

	return autoScalers, nil
}

func SelectAutoScalers(selectors map[string]string) ([]object.AutoScaler, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/autoScalers"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var autoScalers []object.AutoScaler
	err = json.Unmarshal(body, &autoScalers)
	if err != nil {
		log.Println("fail to parse AutoScalers")
		return nil, err
	}

	return autoScalers, nil
}

func CreateAutoScaler(as object.AutoScaler) (object.AutoScaler, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/autoScaler"

	body, err := postRequest(url, as)
	if err != nil {
		log.Println("postRequest fail")
		return as, err
	}

	var newAs object.AutoScaler
	err = json.Unmarshal(body, &newAs)
	if err != nil {
		log.Println("fail to parse AutoScaler")
		return as, err
	}

	return newAs, nil
}

func UpdateAutoScaler(as object.AutoScaler) (object.AutoScaler, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/autoScaler/" + as.UID

	body, err := putRequest(url, as)
	if err != nil {
		log.Println("putRequest fail")
		return as, err
	}

	var newAs object.AutoScaler
	err = json.Unmarshal(body, &newAs)
	if err != nil {
		log.Println("fail to parse AutoScaler")
		return as, err
	}

	return newAs, nil
}

func DeleteAutoScaler(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/autoScaler/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
