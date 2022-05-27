package controller

import (
	"Cubernetes/pkg/actionbrain/informer"
	"Cubernetes/pkg/actionbrain/monitor"
	"Cubernetes/pkg/actionbrain/phase"
	"Cubernetes/pkg/actionbrain/policy"
	"Cubernetes/pkg/actionbrain/types"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/health"
	"Cubernetes/pkg/controllermanager/utils"
	"Cubernetes/pkg/object"
	"log"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	actionUpdateWaitTime = time.Second * 30
	actionScaleWaitTime  = time.Second * 60
	statusUpdateTime     = time.Second * 20
)

type ActionController interface {
	Run()
}

type actionController struct {
	actorInformer  informer.ActorInformer
	actionInformer informer.ActionInformer
	monitor        monitor.ActionMonitor

	biglock sync.Mutex
	wg      *sync.WaitGroup
}

func NewActionController(
	actorInformer informer.ActorInformer,
	actionInformer informer.ActionInformer,
	wg *sync.WaitGroup) (ActionController, error) {
	wg.Add(1)
	monitor, err := monitor.NewActionMonitor()
	if err != nil {
		log.Printf("fail to create ActionMonitor: %v\n", err)
		return nil, err
	}
	return &actionController{
		actorInformer:  actorInformer,
		actionInformer: actionInformer,
		monitor:        monitor,
		biglock:        sync.Mutex{},
		wg:             wg,
	}, nil
}

func (ac *actionController) Run() {
	defer ac.monitor.Close()

	go func() {
		for {
			time.Sleep(time.Second * 11)
			ac.updateActionRoutine()
		}
	}()

	go ac.handleRequest()
	ac.syncLoop()
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
	reqChan := ac.monitor.WatchActionEvoke()

	for req := range reqChan {
		ac.biglock.Lock()

		action := ac.actionInformer.GetMatchedAction(req)
		if action == nil {
			log.Printf("action %s not exist in cache!\n", req)
			continue
		}

		if len(action.Status.ToRun)+len(action.Status.Actors) == 0 {
			// create actor immediately if not exist
			if actor, err := crudobj.CreateActor(ac.buildNewActor(action)); err != nil {
				log.Printf("fail to create actor for action %s: %v", req, err)
			} else {
				log.Printf("create actor %s for action %s", actor.Name, req)
				action.Status.ToRun = append(action.Status.ToRun, actor.UID)

				action.Status.LastScaleTime = time.Now()
				action.Status.LastUpdateTime = time.Now()
				action.Status.DesiredReplicas += 1
				if _, err := crudobj.UpdateAction(*action); err != nil {
					log.Printf("fail to update action status: %v", err)
				}
			}
		}

		ac.biglock.Unlock()
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

			if a.Status != nil {
				ac.checkAndUpdateActionStatus(&a)
			}
		}(action)
	}
	wg.Wait()
}

// simplest scale implementaion: no error-handling
func (ac *actionController) checkAndUpdateActionStatus(action *object.Action) {

	if time.Since(action.Status.LastUpdateTime) < actionUpdateWaitTime {
		return
	}

	actors := ac.actorInformer.GetActors(action.Name)

	runnings := make([]string, 0)
	bads := make([]string, 0)
	for idx := range actors {
		if phase.Running(actors[idx].Status.Phase) {
			runnings = append(runnings, actors[idx].UID)
		} else {
			bads = append(bads, actors[idx].UID)
		}
	}

	desired := action.Status.DesiredReplicas
	lastScale := action.Status.LastScaleTime
	if time.Since(lastScale) > actionScaleWaitTime && len(runnings) == desired {
		// ready to scale
		times, err := ac.monitor.QueryRecentEvoke(action.Name, policy.CountRequestPeriod)
		if err != nil {
			log.Printf("fail to query most recent evoke for %s: %v", action.Name, err)
		} else if target, scale := policy.CalculateScale(times, action.Status.ActualReplicas); scale {
			lastScale = time.Now()
			desired = target
		}
	}

	toCreate := desired - len(runnings)
	toRun := make([]string, 0)
	for idx := 0; idx < toCreate; idx += 1 {
		if actor, err := crudobj.CreateActor(ac.buildNewActor(action)); err != nil {
			log.Printf("fail to create actor to APIServer\n")
		} else {
			log.Printf("Action %s add actor %s\n", action.Name, actor.Name)
			toRun = append(toRun, actor.UID)
		}
	}

	toKill := append(action.Status.ToKill, action.Status.ToRun...)
	if toCreate < 0 {
		toKill = append(toKill, runnings[:-toCreate]...)
	}
	toKill = utils.RemoveDuplication(append(toKill, bads...))
	noExist := make([]int, 0)
	for idx, uid := range toKill {
		if err := crudobj.DeleteActor(uid); err != nil {
			log.Printf("fail to delete actor %s from APIServer: %v\n", uid, err)
			if err.Error() == "fail to delete the obj" {
				noExist = append(noExist, idx)
			}
		} else {
			log.Printf("Action %s remove actor from APIServer: %s\n", action.Name, uid)
		}
	}
	toKill = utils.RemoveMultiIndex(toKill, noExist)

	log.Println("desired:  ", desired)
	log.Println("runnings: ", runnings)
	log.Println("toRun:    ", toRun)
	log.Println("toKill:   ", toKill)

	action.Status = &object.ActionStatus{
		LastScaleTime:   lastScale,
		LastUpdateTime:  time.Now(),
		DesiredReplicas: desired,
		ActualReplicas:  len(runnings),

		Actors: runnings,
		ToRun:  toRun,
		ToKill: toKill,
	}

	if _, err := crudobj.UpdateAction(*action); err != nil {
		log.Printf("fail to update action status of %s: %v\n", action.Name, err)
	} else {
		log.Printf("update action status of %s\n", action.Name)
	}
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
			ActionName:    action.Name,
			ScriptFile:    path.Base(action.Spec.ScriptPath),
			InvokeActions: action.Spec.InvokeActions,
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
