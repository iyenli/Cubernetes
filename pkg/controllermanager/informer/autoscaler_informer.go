package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
)

type AutoScalerInformer interface {
	WatchASEvent() <-chan types.AsEvent
	InformAutoScaler(newAs object.AutoScaler, eType watchobj.EventType) error
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

func (i *cmAutoScalerInformer) WatchASEvent() <-chan types.AsEvent {
	newChan := make(chan types.AsEvent)
	i.asEventChans = append(i.asEventChans, newChan)
	return newChan
}

func (i *cmAutoScalerInformer) InformAutoScaler(newAs object.AutoScaler, eType watchobj.EventType) error {
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
