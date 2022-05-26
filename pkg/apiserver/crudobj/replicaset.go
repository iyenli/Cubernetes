package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetReplicaSet(UID string) (object.ReplicaSet, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/replicaSet/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.ReplicaSet{}, err
	}

	var rs object.ReplicaSet
	err = json.Unmarshal(body, &rs)
	if err != nil {
		log.Println("fail to parse ReplicaSet")
		return object.ReplicaSet{}, err
	}

	return rs, nil
}

func GetReplicaSets() ([]object.ReplicaSet, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/replicaSets"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var rsets []object.ReplicaSet
	err = json.Unmarshal(body, &rsets)
	if err != nil {
		log.Println("fail to parse ReplicaSets")
		return nil, err
	}

	return rsets, nil
}

func SelectReplicaSets(selectors map[string]string) ([]object.ReplicaSet, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/replicaSets"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var rsets []object.ReplicaSet
	err = json.Unmarshal(body, &rsets)
	if err != nil {
		log.Println("fail to parse ReplicaSets")
		return nil, err
	}

	return rsets, nil
}

func CreateReplicaSet(rs object.ReplicaSet) (object.ReplicaSet, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/replicaSet"

	body, err := postRequest(url, rs)
	if err != nil {
		log.Println("postRequest fail")
		return rs, err
	}

	var newRs object.ReplicaSet
	err = json.Unmarshal(body, &newRs)
	if err != nil {
		log.Println("fail to parse ReplicaSet")
		return rs, err
	}

	return newRs, nil
}

func UpdateReplicaSet(rs object.ReplicaSet) (object.ReplicaSet, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/replicaSet/" + rs.UID

	body, err := putRequest(url, rs)
	if err != nil {
		log.Println("putRequest fail")
		return rs, err
	}

	var newRs object.ReplicaSet
	err = json.Unmarshal(body, &newRs)
	if err != nil {
		log.Println("fail to parse ReplicaSet")
		return rs, err
	}

	return newRs, nil
}

func DeleteReplicaSet(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/replicaSet/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("deleteRequest fail")
		return err
	}

	return nil
}
