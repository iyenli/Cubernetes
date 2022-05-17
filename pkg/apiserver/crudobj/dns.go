package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetDns(UID string) (object.Dns, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/dns/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Dns{}, err
	}

	var dns object.Dns
	err = json.Unmarshal(body, &dns)
	if err != nil {
		log.Println("fail to parse Dns")
		return object.Dns{}, err
	}

	return dns, nil
}

func GetDnses() ([]object.Dns, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/dnses"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var dnses []object.Dns
	err = json.Unmarshal(body, &dnses)
	if err != nil {
		log.Println("fail to parse Dnses")
		return nil, err
	}

	return dnses, nil
}

func SelectDnses(selectors map[string]string) ([]object.Dns, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/dnses"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var dnses []object.Dns
	err = json.Unmarshal(body, &dnses)
	if err != nil {
		log.Println("fail to parse Dnses")
		return nil, err
	}

	return dnses, nil
}

func CreateDns(dns object.Dns) (object.Dns, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/dns"

	body, err := postRequest(url, dns)
	if err != nil {
		log.Println("postRequest fail")
		return dns, err
	}

	var newDns object.Dns
	err = json.Unmarshal(body, &newDns)
	if err != nil {
		log.Println("fail to parse Dns")
		return dns, err
	}

	return newDns, nil
}

func UpdateDns(dns object.Dns) (object.Dns, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/dns/" + dns.UID

	body, err := putRequest(url, dns)
	if err != nil {
		log.Println("putRequest fail")
		return dns, err
	}

	var newDns object.Dns
	err = json.Unmarshal(body, &newDns)
	if err != nil {
		log.Println("fail to parse Dns")
		return dns, err
	}

	return newDns, nil
}

func DeleteDns(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/dns/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
