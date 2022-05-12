package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"log"
)

type PodInformer interface {
	WatchPodEvent() <-chan types.PodEvent
	InformPod(newPod object.Pod, eType watchobj.EventType) error
	ListPods() []object.Pod
	CloseChan()
}

func NewPodInformer() (PodInformer, error) {
	return &cubePodInformer{
		podEvent: make(chan types.PodEvent),
		podCache: make(map[string]object.Pod),
	}, nil
}

type cubePodInformer struct {
	podEvent chan types.PodEvent
	podCache map[string]object.Pod
}

func (i *cubePodInformer) WatchPodEvent() <-chan types.PodEvent {
	return i.podEvent
}

func (i *cubePodInformer) CloseChan() {
	close(i.podEvent)
}

func (i *cubePodInformer) ListPods() []object.Pod {
	pods := make([]object.Pod, len(i.podCache))
	idx := 0
	for _, pod := range i.podCache {
		pods[idx] = pod
		idx += 1
	}

	return pods
}

func (i *cubePodInformer) InformPod(newPod object.Pod, eType watchobj.EventType) error {
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.podEvent <- types.PodEvent{
				Type: types.PodRemove,
				Pod:  newPod}
		} else {
			log.Printf("pod %s not exist, DELETE do nothing\n", newPod.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		// update podCache anyway
		i.podCache[newPod.UID] = newPod
		if !exist {
			// UID never seen -> create new Pod
			i.podEvent <- types.PodEvent{
				Type: types.PodCreate,
				Pod:  newPod}
		} else {
			// compute pod change: Name / Label / Spec
			if object.ComputeObjectMetaChange(&newPod.ObjectMeta, &oldPod.ObjectMeta) ||
				object.ComputePodSpecChange(&newPod.Spec, &oldPod.Spec) {
				log.Printf("pod %s spec configured\n", newPod.Name)
				i.podEvent <- types.PodEvent{
					Type: types.PodUpdate,
					Pod:  newPod}
			} else {
				log.Printf("pod %s spec not change\n", newPod.Name)
			}
		}
	}

	return nil
}
