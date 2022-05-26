package action_controller

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/health"
	"Cubernetes/pkg/controllermanager/controller/action_controller/faker"
	"Cubernetes/pkg/controllermanager/informer"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	actionUpdateWaitTime = time.Second * 40
	actionScaleWaitTime = time.Second * 60
	statusUpdateTime = time.Second * 20
)

type ActionController interface {
	Run()
}

type actionController struct {
	actorInformer  informer.ActorInformer
	actionInformer informer.ActionInformer
	recorder       faker.Recorder

	biglock sync.Mutex
	wg      sync.WaitGroup
}

func NewActionController(
	actorInformer informer.ActorInformer,
	actionInformer informer.ActionInformer,
	wg sync.WaitGroup) (ActionController, error) {
	wg.Add(1)
	recorder, _ := faker.NewRecorder()
	return &actionController{
		actorInformer:  actorInformer,
		actionInformer: actionInformer,
		recorder:       recorder,
		biglock:        sync.Mutex{},
		wg:             wg,
	}, nil
}

func (ac *actionController) Run() {

}

func (ac *actionController) syncLoop() {
	actorEventChan := ac.actorInformer.WatchActorEvent()
	defer ac.actorInformer.CloseChan(actorEventChan)

	actionEventChan := ac.actionInformer.WatchActionEvent()
	defer ac.actionInformer.CloseChan(actionEventChan)

	ac.wg.Done()

	for {
		select {
		case actorEvent := <-actorEventChan:
			ac.biglock.Lock()
			actor := actorEvent.Actor
			switch actorEvent.Type {
			case types.ActorCreate:
				err := ac.handleActorCreate(&actor)
				if err != nil {
					log.Printf("fail to handle actor create: %v", err)
				}
			case types.ActorRemove:
				err := ac.handleActorRemove(&actor)
				if err != nil {
					log.Printf("fail to handle actor remove: %v", err)
				}
			default:
				log.Fatal("[FATAL] Unknown actorInformer event types: " + actorEvent.Type)
			}
			ac.biglock.Unlock()
		case actionEvent := <-actionEventChan:
			ac.biglock.Lock()
			action := actionEvent.Action
			switch actionEvent.Type {
			case types.ActionCreate:
				err := ac.handleActionCreate(&action)
				if err != nil {
					log.Printf("fail to handle action create: %v", err)
				}
			case types.ActionRemove:
				err := ac.handleActionRemove(&action)
				if err != nil {
					log.Printf("fail to handle action remove: %v", err)
				}
			default:
				log.Fatal("[FATAL] Unknown actorInformer event types: " + actionEvent.Type)
			}
			ac.biglock.Unlock()
		default:
			time.Sleep(time.Second * 2)
		}
	}

}

// handle request for new function: no sleeping
func (ac *actionController) handleRequest() {
	reqChan := ac.recorder.WatchRequest()

	for req := range reqChan {
		ac.biglock.Lock()
		defer ac.biglock.Unlock()

		action := ac.actionInformer.GetMatchedAction(req)
		if action == nil {
			log.Printf("action %s not exist in cache!\n", req)
			continue
		}

		if len(action.Status.ToRun) + len(action.Status.Actors) == 0 {
			// create actor immediately if not exist
			if actor, err := crudobj.CreateActor(ac.buildNewActor(action)); err != nil {
				log.Printf("fail to create actor for action %s: %v", req, err)
			} else {
				log.Printf("create actor %s for action %s", actor.Name, req)
				action.Status.ToRun = append(action.Status.ToRun, actor.UID)

				action.Status.LastScaleTime = time.Now()
				action.Status.LastUpdateTime = time.Now()
				if _, err := crudobj.UpdateAction(*action); err != nil {
					log.Printf("fail to update action status: %v", err)
				}
			}
		}
	}
}

// scale action
func (ac *actionController) updateActionRoutine() {
	ac.biglock.Lock()
	defer ac.biglock.Unlock()

	if !health.CheckApiServerHealth() {
		log.Printf("[FATAL] lost connection with apiserver: not update this time\n")
		return
	}

	actions := ac.actionInformer.ListActions()

	wg := sync.WaitGroup{}
	wg.Add(len(actions))
	for _, action := range actions {
		go func(a object.Action) {
			defer wg.Done()


		}(action)
	}
	wg.Wait()
}

// simplest scale implementaion: no error-handling
func (ac *actionController) checkAndUpdateActionStatus(action *object.Action) {

}

func (ac *actionController) buildNewActor(action *object.Action) object.Actor {

	actor := object.Actor{
		TypeMeta: object.TypeMeta{
			APIVersion: "v1",
			Kind:       "Actor",
		},
		ObjectMeta: object.ObjectMeta{
			Name: actorName(action.Name),
			Labels: map[string]string{
				"cubernetes.action.uid": action.UID},
		},
		Spec: object.ActorSpec{
			ActionName: action.Name,
			ScriptFile: path.Base(action.Spec.ScriptPath),
		},
		Status: &object.ActorStatus{
			Phase: object.ActorCreated,
		},
	}

	return actor
}

func actorName(actionName string) string {
	return strings.Join([]string{actionName, uuid.New().String()[:8]}, "_")
}
