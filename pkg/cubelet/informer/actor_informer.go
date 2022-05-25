package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type ActorInformer interface {
	ListAndWatchActorsWithRetry()
	WatchActorEvent() <-chan types.ActorEvent
	SetNodeUID(uid string)
	CloseChan()
}

func NewActorInformer() (ActorInformer, error) {
	return &cubeActorInformer{
		actorEventChan: make(chan types.ActorEvent),
		actorCache:     make(map[string]object.Actor),
	}, nil
}

type cubeActorInformer struct {
	actorEventChan chan types.ActorEvent
	actorCache     map[string]object.Actor
	nodeUID        string
}

func (c *cubeActorInformer) SetNodeUID(uid string) {
	if c.nodeUID != "" {
		log.Printf("[FATAL]: Node ID already set!\n")
	}
	c.nodeUID = uid
}

func (c *cubeActorInformer) WatchActorEvent() <-chan types.ActorEvent {
	return c.actorEventChan
}

func (c *cubeActorInformer) CloseChan() {
	close(c.actorEventChan)
}

func (c *cubeActorInformer) ListAndWatchActorsWithRetry() {
	defer close(c.actorEventChan)
	for {
		c.tryListandWatchActors()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (c *cubeActorInformer) tryListandWatchActors() {
	if allActors, err := crudobj.GetActors(); err != nil {
		log.Printf("[INFO]: fail to get all actors from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, actor := range allActors {
			if actor.Status != nil && actor.Status.NodeUID == c.nodeUID &&
				actor.Status.Phase == object.ActorCreated {
				c.informActor(actor, watchobj.EVENT_PUT)
			}
		}
	}

	ch, cancel, err := watchobj.WatchActors()
	if err != nil {
		log.Printf("fail to watch actors from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case actorEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			}
			if actorEvent.Actor.Status == nil && actorEvent.EType != watchobj.EVENT_DELETE {
				log.Println("[INFO]: Actor caught, but status is nil so Cubelet doesn't handle it")
				continue
			}
			if actorEvent.EType == watchobj.EVENT_DELETE || actorEvent.Actor.Status.NodeUID == c.nodeUID {
				log.Println("[INFO]: my actor caught, types is", actorEvent.EType)
				switch actorEvent.EType {
				case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
					c.informActor(actorEvent.Actor, actorEvent.EType)
				default:
					log.Panic("[Error]: Unsupported types in watch actor.")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (c *cubeActorInformer) informActor(actor object.Actor, eType watchobj.EventType) {
	_, exist := c.actorCache[actor.UID]
	if eType == watchobj.EVENT_DELETE && exist {
		c.actorEventChan <- types.ActorEvent{
			Type:  types.Remove,
			Actor: actor,
		}
	} else if eType == watchobj.EVENT_PUT && !exist {
		c.actorEventChan <- types.ActorEvent{
			Type:  types.Create,
			Actor: actor,
		}
	} else if !exist {
		log.Printf("[Error] put actor %s not exist\n", actor.Name)
	}
}
