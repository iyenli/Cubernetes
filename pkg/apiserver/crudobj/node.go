package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetNode(UID string) (object.Node, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/node/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Node{}, err
	}

	var node object.Node
	err = json.Unmarshal(body, &node)
	if err != nil {
		log.Println("fail to parse Node")
		return object.Node{}, err
	}

	return node, nil
}

func GetNodes() ([]object.Node, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/nodes"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var nodes []object.Node
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Println("fail to parse Nodes")
		return nil, err
	}

	return nodes, nil
}

func SelectNodes(selectors map[string]string) ([]object.Node, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/nodes"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var nodes []object.Node
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		log.Println("fail to parse Nodes")
		return nil, err
	}

	return nodes, nil
}

func CreateNode(node object.Node) (object.Node, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/node"

	body, err := postRequest(url, node)
	if err != nil {
		log.Println("postRequest fail")
		return node, err
	}

	var newNode object.Node
	err = json.Unmarshal(body, &newNode)
	if err != nil {
		log.Println("fail to parse Node, body: ", string(body))
		return node, err
	}

	return newNode, nil
}

func UpdateNode(node object.Node) (object.Node, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/node/" + node.UID

	body, err := putRequest(url, node)
	if err != nil {
		log.Println("putRequest fail")
		return node, err
	}

	var newNode object.Node
	err = json.Unmarshal(body, &newNode)
	if err != nil {
		log.Println("fail to parse Node")
		return node, err
	}

	return newNode, nil
}

func DeleteNode(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/node/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
