package action_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"log"
)

func (ac *actionController) handleActionCreate(action *object.Action) error {
	// copy script to apiserver

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
