package scheduler

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/scheduler/types"
	"log"
	"time"
)

func (sr *ScheduleRuntime) ScheduleActor(Actor *object.Actor) {
	if Actor.Status == nil {
		log.Println("[INFO]: ignore actor whose status is nil")
		return
	}

	if Actor.Status.Phase == object.ActorCreated && Actor.Status.NodeUID == "" {
		ActorInfo, err := sr.Implement.Schedule()
		if err != nil {
			log.Println("[Error]: when scheduling, error:", err.Error())
			return
		}

		if ActorInfo.NodeUUID == "" {
			log.Println("[Warn]: No node to schedule this node")
		}

		err = sr.SendActorScheduleInfoBack(Actor, &ActorInfo)
		if err != nil {
			log.Println("[Error]: when sending scheduler result,", err.Error())
			return
		}
	}
}

func (sr *ScheduleRuntime) SendActorScheduleInfoBack(ActorToSchedule *object.Actor, info *types.ScheduleInfo) error {
	ActorToSchedule.Status.NodeUID = info.NodeUUID
	ActorToSchedule.Status.Phase = object.ActorBound

	_, err := crudobj.UpdateActor(*ActorToSchedule)
	if err != nil {
		log.Println("[INFO]: Update Actor failed")
		return err
	}

	log.Println("[INFO]: Schedule Actor", ActorToSchedule.UID, "to node", info.NodeUUID)
	return nil
}

func (sr *ScheduleRuntime) WatchActor() {
	for true {
		sr.tryWatchActor()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (sr *ScheduleRuntime) tryWatchActor() {
	if allActors, err := crudobj.GetActors(); err != nil {
		log.Printf("[INFO]: fail to get all Actors from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, Actor := range allActors {
			sr.ScheduleActor(&Actor)
		}
	}

	ch, cancel, err := watchobj.WatchActors()
	if err != nil {
		log.Printf("[Error]: Error occurs when watching Actors: %v", err)
		return
	}
	defer cancel()

	for true {
		select {
		case ActorEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			} else {
				switch ActorEvent.EType {
				case watchobj.EVENT_PUT:
					sr.ScheduleActor(&ActorEvent.Actor)
				case watchobj.EVENT_DELETE:
					log.Println("[Info]: Delete Actor, do nothing")
				default:
					log.Panic("[Fatal]: Unsupported types in watching Actor.")
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
