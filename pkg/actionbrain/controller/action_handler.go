package controller

import (
	"Cubernetes/pkg/actionbrain/phase"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/kafka"
	"log"
	"time"
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

func (ac *actionController) handleActionUpdate(action *object.Action) error {
	// only handle script change
	actors := ac.actorInformer.GetActors(action.Name)
	for _, actor := range actors {
		if phase.Running(actor.Status.Phase) {
			actor.Spec.ScriptUID = action.Spec.ScriptUID
			actor.Status.LastUpdatedTime = time.Now()
			if _, err := crudobj.UpdateActor(actor); err != nil {
				log.Printf("fail to update script for Actot %s\n", actor.Name)
			}
		}
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
