package controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/kafka"
	"log"
)

func (ac *actionController) handleActionCreate(action *object.Action) error {
	// create action topic
	topicName := action.Name + "_TOPIC"
	if err := kafka.CreateTopic(ac.kafkaHost, topicName); err != nil {
		log.Printf("fail to create receive-topic for Action %s\n", action.Name)
		return err
	}

	return nil
}

func (ac *actionController) handleActionRemove(action *object.Action) error {
	for _, uid := range action.Status.Actors {
		if err := crudobj.DeleteActor(uid); err != nil {
			log.Printf("fail to delete actor %s from API Server: %v\n", uid, err)
		} else {
			log.Printf("Action %s remove actor from API Server: %s\n", action.Name, uid)
		}
	}

	return nil
}
