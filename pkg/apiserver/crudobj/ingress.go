package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetIngress(UID string) (object.Ingress, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/ingress/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Ingress{}, err
	}

	var ingress object.Ingress
	err = json.Unmarshal(body, &ingress)
	if err != nil {
		log.Println("fail to parse Ingress")
		return object.Ingress{}, err
	}

	return ingress, nil
}

func GetIngresses() ([]object.Ingress, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/ingresses"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var ingresses []object.Ingress
	err = json.Unmarshal(body, &ingresses)
	if err != nil {
		log.Println("fail to parse Ingresses")
		return nil, err
	}

	return ingresses, nil
}

func SelectIngresses(selectors map[string]string) ([]object.Ingress, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/ingresses"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var ingresses []object.Ingress
	err = json.Unmarshal(body, &ingresses)
	if err != nil {
		log.Println("fail to parse Ingresses")
		return nil, err
	}

	return ingresses, nil
}

func CreateIngress(ingress object.Ingress) (object.Ingress, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/ingress"

	body, err := postRequest(url, ingress)
	if err != nil {
		log.Println("postRequest fail")
		return ingress, err
	}

	var newIngress object.Ingress
	err = json.Unmarshal(body, &newIngress)
	if err != nil {
		log.Println("fail to parse Ingress")
		return ingress, err
	}

	return newIngress, nil
}

func UpdateIngress(ingress object.Ingress) (object.Ingress, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/ingress/" + ingress.UID

	body, err := putRequest(url, ingress)
	if err != nil {
		log.Println("putRequest fail")
		return ingress, err
	}

	var newIngress object.Ingress
	err = json.Unmarshal(body, &newIngress)
	if err != nil {
		log.Println("fail to parse Ingress")
		return ingress, err
	}

	return newIngress, nil
}

func DeleteIngress(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/ingress/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("deleteRequest fail")
		return err
	}

	return nil
}
