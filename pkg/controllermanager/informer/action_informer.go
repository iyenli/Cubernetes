package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type ActionInformer interface {
	ListAndWatchActionsWithRetry()
	WatchActionEvent() <-chan types.ActionEvent
	GetMatchedAction(actionName string) *object.Action
	ListActions() []object.Action
	CloseChan(<-chan types.ActionEvent)
}

const watchActionRetryIntervalSec = 16

func NewActionInformer() (ActionInformer, error) {
	return &cmActionInformer{
		actionEventChans: make([]chan types.ActionEvent, 0),
		actionCache:      make(map[string]object.Action),
	}, nil
}

type cmActionInformer struct {
	actionEventChans []chan types.ActionEvent
	actionCache      map[string]object.Action
}

func (i *cmActionInformer) WatchActionEvent() <-chan types.ActionEvent {
	newChan := make(chan types.ActionEvent)
	i.actionEventChans = append(i.actionEventChans, newChan)
	return newChan
}

func (i *cmActionInformer) CloseChan(ch <-chan types.ActionEvent) {
	found := -1
	for idx, c := range i.actionEventChans {
		if c == ch {
			close(c)
			found = idx
			break
		}
	}
	if found != -1 {
		i.actionEventChans = append(i.actionEventChans[:found], i.actionEventChans[found+1:]...)
	}
}

func (i *cmActionInformer) ListAndWatchActionsWithRetry() {
	for {
		i.tryListAndWatchActions()
		time.Sleep(watchActionRetryIntervalSec * time.Second)
	}
}

func (i *cmActionInformer) GetMatchedAction(actionName string) *object.Action {
	for _, action := range i.actionCache {
		if action.Name == actionName {
			return &action
		}
	}
	return nil
}

func (i *cmActionInformer) ListActions() []object.Action {
	actions := make([]object.Action, 0)
	for _, action := range i.actionCache {
		actions = append(actions, action)
	}
	return actions
}

func (i *cmActionInformer) tryListAndWatchActions() {
	if allActions, err := crudobj.GetActions(); err != nil {
		return
	} else {
		for _, action := range allActions {
			i.actionCache[action.UID] = action
		}
	}

	ch, cancel, err := watchobj.WatchActions()
	if err != nil {
		log.Printf("fail to watch action from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case actionEvent, ok := <-ch:
			if !ok {
				log.Printf("lost connection with APIServer, retry after %d seconds...\n", watchActionRetryIntervalSec)
				return
			}
			action := actionEvent.Action
			switch actionEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				i.informAction(action, actionEvent.EType)
			default:
				log.Fatal("[FATAL] Unknown event types: " + actionEvent.EType)
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *cmActionInformer) informAction(newAction object.Action, eType watchobj.EventType) {
	oldAction, exist := i.actionCache[newAction.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.actionCache, newAction.UID)
			i.informAll(types.ActionEvent{
				Type:   types.ActionRemove,
				Action: oldAction,
			})
		} else {
			log.Printf("action %s not exist, DELETE do nothing\n", newAction.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		if !exist {
			i.actionCache[newAction.UID] = newAction
			i.informAll(types.ActionEvent{
				Type:   types.ActionCreate,
				Action: newAction,
			})
		} else {
			log.Printf("[Error] update action not supported!\n")
		}
	}
}

func (i *cmActionInformer) informAll(event types.ActionEvent) {
	for _, c := range i.actionEventChans {
		c <- event
	}
}
