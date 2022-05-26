package controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"fmt"
	"log"
	"time"
)

func (ac *actionController) handleActorCreate(actor *object.Actor) error {
	action := ac.actionInformer.GetMatchedAction(actor.Spec.ActionName)
	if action == nil {
		return fmt.Errorf("action %s not found in cache", actor.Spec.ActionName)
	}

	if _, found := ac.actorUIDAppearedIndex(actor.UID, action.Status.Actors); !found {
		action.Status.ActualReplicas += 1
		action.Status.Actors = append(action.Status.Actors, actor.UID)
	}

	if idx, found := ac.actorUIDAppearedIndex(actor.UID, action.Status.ToRun); found {
		action.Status.ToRun =
			append(action.Status.ToRun[:idx], action.Status.ToRun[idx+1:]...)
	} else {
		log.Printf("[FATAL] unexpected actor %s add to Action %s when create\n", actor.Name, action.Name)
	}

	action.Status.LastUpdateTime = time.Now()
	if _, err := crudobj.UpdateAction(*action); err != nil {
		log.Printf("fail to update Action status to apiserver\n")
		return err
	}

	return nil
}

func (ac *actionController) handleActorRemove(actor *object.Actor) error {
	action := ac.actionInformer.GetMatchedAction(actor.Spec.ActionName)
	if action == nil {
		return fmt.Errorf("action %s not found in cache", actor.Spec.ActionName)
	}

	if idx, found := ac.actorUIDAppearedIndex(actor.UID, action.Status.ToKill); found {
		action.Status.ToKill =
			append(action.Status.ToKill[:idx], action.Status.ToKill[idx+1:]...)
	}

	if idx, found := ac.actorUIDAppearedIndex(actor.UID, action.Status.Actors); found {
		action.Status.ActualReplicas -= 1
		action.Status.Actors =
			append(action.Status.Actors[:idx], action.Status.Actors[idx+1:]...)
	} else {
		log.Printf("actor %s killed but not in running\n", actor.Name)
	}

	if _, err := crudobj.UpdateAction(*action); err != nil {
		log.Printf("fail to update Action status to apiserver: %v\n", err)
		return err
	}

	return nil
}

func (ac *actionController) actorUIDAppearedIndex(actorUID string, list []string) (int, bool) {
	for idx, uid := range list {
		if uid == actorUID {
			return idx, true
		}
	}
	return -1, false
}
