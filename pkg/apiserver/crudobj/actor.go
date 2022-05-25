package crudobj

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"log"
	"strconv"
)

func GetActor(UID string) (object.Actor, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actor/" + UID

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return object.Actor{}, err
	}

	var actor object.Actor
	err = json.Unmarshal(body, &actor)
	if err != nil {
		log.Println("fail to parse Actor")
		return object.Actor{}, err
	}

	return actor, nil
}

func GetActors() ([]object.Actor, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actors"

	body, err := getRequest(url)
	if err != nil {
		log.Println("getRequest fail")
		return nil, err
	}

	var actors []object.Actor
	err = json.Unmarshal(body, &actors)
	if err != nil {
		log.Println("fail to parse Actors")
		return nil, err
	}

	return actors, nil
}

func SelectActors(selectors map[string]string) ([]object.Actor, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/select/actors"

	body, err := postRequest(url, selectors)
	if err != nil {
		log.Println("postRequest fail")
		return nil, err
	}

	var actors []object.Actor
	err = json.Unmarshal(body, &actors)
	if err != nil {
		log.Println("fail to parse Actors")
		return nil, err
	}

	return actors, nil
}

func CreateActor(actor object.Actor) (object.Actor, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actor"

	body, err := postRequest(url, actor)
	if err != nil {
		log.Println("postRequest fail")
		return actor, err
	}

	var newActor object.Actor
	err = json.Unmarshal(body, &newActor)
	if err != nil {
		log.Println("fail to parse Actor")
		return actor, err
	}

	return newActor, nil
}

func UpdateActor(actor object.Actor) (object.Actor, error) {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actor/" + actor.UID

	body, err := putRequest(url, actor)
	if err != nil {
		log.Println("putRequest fail")
		return actor, err
	}

	var newActor object.Actor
	err = json.Unmarshal(body, &newActor)
	if err != nil {
		log.Println("fail to parse Actor")
		return actor, err
	}

	return newActor, nil
}

func DeleteActor(UID string) error {
	url := "http://" + cubeconfig.APIServerIp + ":" + strconv.Itoa(cubeconfig.APIServerPort) + "/apis/actor/" + UID

	err := deleteRequest(url)
	if err != nil {
		log.Println("postRequest fail")
		return err
	}

	return nil
}
