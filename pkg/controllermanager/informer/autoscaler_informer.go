package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type AutoScalerInformer interface {
	ListAndWatchAutoScalersWithRetry()
	WatchASEvent() <-chan types.AsEvent
	ListAutoScalers() []object.AutoScaler
	GetAutoScaler(UID string) (*object.AutoScaler, bool)
	CloseChan(<-chan types.AsEvent)
}

func NewAutoScalerInformer() (AutoScalerInformer, error) {
	return &cmAutoScalerInformer{
		asEventChans: make([]chan types.AsEvent, 0),
		asCache:      make(map[string]object.AutoScaler),
	}, nil
}

type cmAutoScalerInformer struct {
	asEventChans []chan types.AsEvent
	asCache      map[string]object.AutoScaler
}

func (i *cmAutoScalerInformer) ListAndWatchAutoScalersWithRetry() {
	for {
		i.tryListAndWatchAutoScalers()
		time.Sleep(watchASRetryIntervalSec * time.Second)
	}
}

func (i *cmAutoScalerInformer) tryListAndWatchAutoScalers() {
	if all, err := crudobj.GetAutoScalers(); err != nil {
		log.Printf("[Manager] fail to get all AutoScalers from apiserver: %v\n", err)
		log.Printf("[Manager] will retry after %d seconds...\n", watchASRetryIntervalSec)
		return
	} else {
		for _, as := range all {
			i.asCache[as.UID] = as
		}
	}

	ch, cancel, err := watchobj.WatchAutoScalers()
	if err != nil {
		log.Printf("fail to watch AutoScalers from apiserver: %v\n", err)
		return
	}
	defer cancel()

	for {
		select {
		case asEvent, ok := <-ch:
			if !ok {
				log.Printf("lost connection with APIServer, retry after %d seconds...\n", watchASRetryIntervalSec)
				return
			}
			as := asEvent.AutoScaler
			switch asEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				i.informAutoScaler(as, asEvent.EType)
			default:
				log.Fatal("[FATAL] Unknown event types: " + asEvent.EType)
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *cmAutoScalerInformer) WatchASEvent() <-chan types.AsEvent {
	newChan := make(chan types.AsEvent, 10)
	i.asEventChans = append(i.asEventChans, newChan)
	return newChan
}

func (i *cmAutoScalerInformer) informAutoScaler(newAs object.AutoScaler, eType watchobj.EventType) error {
	oldAs, exist := i.asCache[newAs.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.asCache, newAs.UID)
			i.informAll(types.AsEvent{
				Type:       types.AsRemove,
				AutoScaler: oldAs,
			})
		} else {
			log.Printf("AutoScaler %s not exist, DELETE do nothing\n", newAs.UID)
		}
	}

	if eType == watchobj.EVENT_PUT {
		if !exist {
			i.asCache[newAs.UID] = newAs
			i.informAll(types.AsEvent{
				Type:       types.AsCreate,
				AutoScaler: newAs,
			})
		} else {
			if object.ComputeAutoScalerSpecChange(&newAs.Spec, &oldAs.Spec) {
				log.Printf("[FATAL] AutoScaler Spec change is not supported!\n")
			} else {
				i.asCache[newAs.UID] = newAs
				i.informAll(types.AsEvent{
					Type:       types.AsUpdate,
					AutoScaler: newAs,
				})
			}
		}
	}

	return nil
}

func (i *cmAutoScalerInformer) ListAutoScalers() []object.AutoScaler {
	autoScalers := make([]object.AutoScaler, len(i.asCache))
	idx := 0
	for _, as := range i.asCache {
		autoScalers[idx] = as
		idx += 1
	}

	return autoScalers
}

func (i *cmAutoScalerInformer) GetAutoScaler(UID string) (*object.AutoScaler, bool) {
	as, ok := i.asCache[UID]
	if ok {
		return &as, true
	} else {
		return nil, false
	}
}

func (i *cmAutoScalerInformer) CloseChan(ch <-chan types.AsEvent) {
	found := -1
	for idx, c := range i.asEventChans {
		if c == ch {
			close(c)
			found = idx
			break
		}
	}
	if found != -1 {
		i.asEventChans = append(i.asEventChans[:found], i.asEventChans[found+1:]...)
	}
}

func (i *cmAutoScalerInformer) informAll(event types.AsEvent) {
	for _, c := range i.asEventChans {
		c <- event
	}
}
