package replicaset_controller

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/object"
	"log"
)

type ReplicaSetInformer interface {
	WatchRSEvent() <-chan rsEvent
	InformReplicaSet(newRs object.ReplicaSet, eType watchobj.EventType) error
	GetMatchedReplicaSet(pod *object.Pod) []object.ReplicaSet
	ListReplicaSets() []object.ReplicaSet
	CloseChan()
}

func NewReplicaSetInformer() (ReplicaSetInformer, error) {
	return &rsInformer{
		rsEvent: make(chan rsEvent),
		rsCache: make(map[string]object.ReplicaSet),
	}, nil
}

type rsInformer struct {
	rsEvent chan rsEvent
	rsCache map[string]object.ReplicaSet
}

type rsEventType string

const (
	rsCreate rsEventType = "create"
	rsUpdate rsEventType = "update"
	rsRemove rsEventType = "remove"
)

type rsEvent struct {
	Type       rsEventType
	ReplicaSet object.ReplicaSet
}

func (i *rsInformer) WatchRSEvent() <-chan rsEvent {
	return i.rsEvent
}

func (i *rsInformer) InformReplicaSet(newRs object.ReplicaSet, eType watchobj.EventType) error {
	oldRs, exist := i.rsCache[newRs.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			// delete ReplicaSet from cache
			// and tell apiserver to delete its pod
			delete(i.rsCache, newRs.UID)
			i.rsEvent <- rsEvent{
				Type:       rsRemove,
				ReplicaSet: oldRs,
			}
		} else {
			log.Printf("ReplicaSet %s not exist, DELETE do nothing\n", newRs.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		i.rsCache[newRs.UID] = newRs
		if !exist {
			i.rsEvent <- rsEvent{
				Type:       rsCreate,
				ReplicaSet: newRs,
			}
		} else {
			if newRs.Spec.Replicas != oldRs.Spec.Replicas {
				log.Printf("ReplicaSet %s spec replicas configured\n", newRs.Name)
				i.rsEvent <- rsEvent{
					Type:       rsUpdate,
					ReplicaSet: newRs,
				}
			} else if object.ComputeReplicaSetSpecChange(&newRs.Spec, &oldRs.Spec) {
				log.Println("[UNHANDLED] PodSpec of ReplicaSet update is not handled")
			} else {
				log.Printf("ReplicaSet %s spec not change\n", newRs.Name)
			}
		}
	}

	return nil
}

func (i *rsInformer) GetMatchedReplicaSet(pod *object.Pod) []object.ReplicaSet {
	matched := make([]object.ReplicaSet, 0)
	for _, rs := range i.rsCache {
		if object.MatchLabelSelector(rs.Spec.Selector, pod.Labels) {
			matched = append(matched, rs)
		}
	}
	return matched
}

func (i *rsInformer) ListReplicaSets() []object.ReplicaSet {
	replicaSets := make([]object.ReplicaSet, len(i.rsCache))
	idx := 0
	for _, rs := range i.rsCache {
		replicaSets[idx] = rs
		idx += 1
	}

	return replicaSets
}

func (i *rsInformer) CloseChan() {
	close(i.rsEvent)
}
