package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type ReplicaSetInformer interface {
	ListAndWatchReplicaSetsWithRetry()
	WatchRSEvent() <-chan types.RsEvent
	GetMatchedReplicaSet(pod *object.Pod) []object.ReplicaSet
	ListReplicaSets() []object.ReplicaSet
	GetReplicaSet(UID string) (*object.ReplicaSet, bool)
	CloseChan(<-chan types.RsEvent)
}

func NewReplicaSetInformer() (ReplicaSetInformer, error) {
	return &rsInformer{
		RsEventChans: make([]chan types.RsEvent, 0),
		rsCache:      make(map[string]object.ReplicaSet),
	}, nil
}

type rsInformer struct {
	RsEventChans []chan types.RsEvent
	rsCache      map[string]object.ReplicaSet
}

func (i *rsInformer) ListAndWatchReplicaSetsWithRetry() {
	for {
		i.tryListAndWatchReplicaSets()
		time.Sleep(watchRSRetryIntervalSec * time.Second)
	}
}

func (i *rsInformer) tryListAndWatchReplicaSets() {

	if all, err := crudobj.GetReplicaSets(); err != nil {
		log.Printf("[Manager] fail to get all ReplicaSets from apiserver: %v\n", err)
		log.Printf("[Manager] will retry after %d seconds...\n", watchRSRetryIntervalSec)
		return
	} else {
		for _, rs := range all {
			i.rsCache[rs.UID] = rs
		}
	}

	ch, cancel, err := watchobj.WatchReplicaSets()
	if err != nil {
		log.Printf("fail to watch ReplicaSets from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case rsEvent, ok := <-ch:
			if !ok {
				log.Printf("lost connection with APIServer, retry after %d seconds...\n", watchRSRetryIntervalSec)
				return
			}
			rs := rsEvent.ReplicaSet
			switch rsEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				i.informReplicaSet(rs, rsEvent.EType)
			default:
				log.Fatal("[FATAL] Unknown event types: " + rsEvent.EType)
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *rsInformer) WatchRSEvent() <-chan types.RsEvent {
	log.Printf("replicaset informer make a new chan!\n")
	newChan := make(chan types.RsEvent)
	i.RsEventChans = append(i.RsEventChans, newChan)
	return newChan
}

func (i *rsInformer) informReplicaSet(newRs object.ReplicaSet, eType watchobj.EventType) error {
	oldRs, exist := i.rsCache[newRs.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			// delete ReplicaSet from cache
			// and tell apiserver to delete its pod
			delete(i.rsCache, newRs.UID)
			i.informAll(types.RsEvent{
				Type:       types.RsRemove,
				ReplicaSet: oldRs,
			})
		} else {
			log.Printf("ReplicaSet %s not exist, DELETE do nothing\n", newRs.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		i.rsCache[newRs.UID] = newRs
		if !exist {
			i.informAll(types.RsEvent{
				Type:       types.RsCreate,
				ReplicaSet: newRs,
			})
		} else {
			if newRs.Spec.Replicas != oldRs.Spec.Replicas {
				log.Printf("ReplicaSet %s spec replicas configured\n", newRs.Name)
				i.informAll(types.RsEvent{
					Type:       types.RsUpdate,
					ReplicaSet: newRs,
				})
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

func (i *rsInformer) GetReplicaSet(UID string) (*object.ReplicaSet, bool) {
	rs, ok := i.rsCache[UID]
	if ok {
		return &rs, true
	} else {
		return nil, false
	}
}

func (i *rsInformer) CloseChan(ch <-chan types.RsEvent) {
	log.Printf("replicaset informer close a chan!\n")
	found := -1
	for idx, c := range i.RsEventChans {
		if c == ch {
			close(c)
			found = idx
			break
		}
	}
	if found != -1 {
		i.RsEventChans = append(i.RsEventChans[:found], i.RsEventChans[found+1:]...)
	}
}

func (i *rsInformer) informAll(event types.RsEvent) {
	for _, c := range i.RsEventChans {
		c <- event
	}
}
