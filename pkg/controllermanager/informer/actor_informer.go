package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/phase"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type ActorInformer interface {
	ListAndWatchActorsWithRetry()
	WatchActorEvent() <-chan types.ActorEvent
	GetActors(actionName string) []object.Actor
	CloseChan(<-chan types.ActorEvent)
}

const (
	watchActorsRetryIntervalSec = 16
)

func NewActorInformer() (ActorInformer, error) {
	return &cmActorInformer{
		actorEventChans: make([]chan types.ActorEvent, 0),
		actorCache:      make(map[string]object.Actor),
	}, nil
}

type cmActorInformer struct {
	actorEventChans []chan types.ActorEvent
	actorCache      map[string]object.Actor
}

func (i *cmActorInformer) WatchActorEvent() <-chan types.ActorEvent {
	newChan := make(chan types.ActorEvent)
	i.actorEventChans = append(i.actorEventChans, newChan)
	return newChan
}

func (i *cmActorInformer) CloseChan(ch <-chan types.ActorEvent) {
	found := -1
	for idx, c := range i.actorEventChans {
		if c == ch {
			close(c)
			found = idx
			break
		}
	}
	if found != -1 {
		i.actorEventChans = append(i.actorEventChans[:found], i.actorEventChans[found+1:]...)
	}
}

func (i *cmActorInformer) ListAndWatchActorsWithRetry() {
	for {
		i.tryListAndWatchActors()
		time.Sleep(watchActorsRetryIntervalSec * time.Second)
	}
}

func (i *cmActorInformer) GetActors(actionName string) []object.Actor {
	actors := make([]object.Actor, 0)
	for _, actor := range i.actorCache {
		if actor.Spec.ActionName == actionName {
			actors = append(actors, actor)
		}
	}
	return actors
}

func (i *cmActorInformer) tryListAndWatchActors() {
	if allActors, err := crudobj.GetActors(); err != nil {
		log.Printf("[Manager] fail to get all actors from apiserver: %v\n", err)
		log.Printf("[Manager] will retry after %d seconds...\n", watchActorsRetryIntervalSec)
		return
	} else {
		for _, actor := range allActors {
			if actor.Status != nil && !phase.ActorNotHandle(actor.Status.Phase) {
				i.actorCache[actor.UID] = actor
			}
		}
	}

	ch, cancel, err := watchobj.WatchActors()
	if err != nil {
		log.Printf("fail to watch actor from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case actorEvent, ok := <-ch:
			if !ok {
				log.Printf("lost connection with APIServer, retry after %d seconds...\n", watchActorsRetryIntervalSec)
				return
			}
			actor := actorEvent.Actor
			if (actor.Status == nil || phase.ActorNotHandle(actor.Status.Phase)) &&
				actorEvent.EType != watchobj.EVENT_DELETE {
				continue
			}
			switch actorEvent.EType {
			case watchobj.EVENT_DELETE, watchobj.EVENT_PUT:
				err := i.informActor(actor, actorEvent.EType)
				if err != nil {
					log.Println("[INFO]: Delete or put an actor error", actor.UID)
					return
				}
			default:
				log.Panic("Unsupported types in watch actor.")
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *cmActorInformer) informActor(newActor object.Actor, eType watchobj.EventType) error {
	oldActor, exist := i.actorCache[newActor.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.actorCache, newActor.UID)
			i.informAll(types.ActorEvent{
				Type:  types.ActorRemove,
				Actor: oldActor,
			})
		} else {
			log.Printf("actor %s not exist, DELETE do nothing\n", newActor.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		if !exist && phase.ActorRunning(newActor.Status.Phase) {
			i.actorCache[newActor.UID] = newActor
			i.informAll(types.ActorEvent{
				Type:  types.ActorCreate,
				Actor: newActor,
			})
		} else if exist && !object.ComputeActorSpecChange(&newActor, &oldActor) {
			// just update status
			i.actorCache[newActor.UID] = newActor
		} else {
			log.Printf("[Error] something bad happend!\n")
		}
	}

	return nil
}

func (i *cmActorInformer) informAll(event types.ActorEvent) {
	for _, c := range i.actorEventChans {
		c <- event
	}
}
