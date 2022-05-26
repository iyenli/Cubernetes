package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetAction(UID string) (object.Action, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Action{}, err
	}

	var action object.Action
	err = json.Unmarshal(body, &action)
	if err != nil {
		log.Println("fail to parse Action")
		return object.Action{}, err
	}

	return action, nil
}

func GetActions() ([]object.Action, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actions"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var actions []object.Action
	err = json.Unmarshal(body, &actions)
	if err != nil {
		log.Println("fail to parse Actions")
		return nil, err
	}

	return actions, nil
}

func SelectActions(selectors map[string]string) ([]object.Action, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/actions"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var actions []object.Action
	err = json.Unmarshal(body, &actions)
	if err != nil {
		log.Println("fail to parse Actions")
		return nil, err
	}

	return actions, nil
}

func CreateAction(action object.Action) (object.Action, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action"

	body, err := postRequest(url, action)
	if err != nil {
		log.Println("postRequest fail")
		return action, err
	}

	var newAction object.Action
	err = json.Unmarshal(body, &newAction)
	if err != nil {
		log.Println("fail to parse Action")
		return action, err
	}

	return newAction, nil
}

func UpdateAction(action object.Action) (object.Action, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/" + action.UID

	body, err := putRequest(url, action)
	if err != nil {
		log.Println("putRequest fail")
		return action, err
	}

	var newAction object.Action
	err = json.Unmarshal(body, &newAction)
	if err != nil {
		log.Println("fail to parse Action")
		return action, err
	}

	return newAction, nil
}

func DeleteAction(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/action/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("deleteRequest fail")
		return err
	}

	return nil
}
