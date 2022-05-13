package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
)

type PodInformer interface {
	WatchPodEvent() <-chan types.PodEvent
	InformPod(newPod object.Pod, eType watchobj.EventType) error
	ListPods() []object.Pod
	CloseChan()
}

type ProxyPodInformer struct {
	podChannel chan types.PodEvent
	podCache   map[string]object.Pod
}

func NewPodInformer() PodInformer {
	return &ProxyPodInformer{
		podChannel: make(chan types.PodEvent),
		podCache:   make(map[string]object.Pod),
	}
}

func (i *ProxyPodInformer) WatchPodEvent() <-chan types.PodEvent {
	return i.podChannel
}

func (i *ProxyPodInformer) CloseChan() {
	close(i.podChannel)
}

func (i *ProxyPodInformer) ListPods() []object.Pod {
	pods := make([]object.Pod, len(i.podCache))
	idx := 0
	for _, pod := range i.podCache {
		pods[idx] = pod
		idx += 1
	}

	return pods
}

func (i *ProxyPodInformer) InformPod(newPod object.Pod, eType watchobj.EventType) error {
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.podChannel <- types.PodEvent{
				Type: types.PodRemove,
				Pod:  newPod}
		} else {
			log.Printf("pod %s not exist, delete do nothing\n", newPod.UID)
		}
	}

	if eType == watchobj.EVENT_PUT {
		// Just handle pod whose ip has been allocated
		if newPod.Status.IP == nil {
			log.Println("Pod without ip, just ignore")
			return nil
		}

		// else cached and judge type
		i.podCache[newPod.UID] = newPod
		if !exist {
			i.podChannel <- types.PodEvent{
				Type: types.PodCreate,
				Pod:  newPod}
		} else {
			// compute pod change: IP / Label
			if object.ComputePodNetworkChange(&newPod, &oldPod) {
				log.Println("pod changed, pod ID is:", newPod.UID)
				i.podChannel <- types.PodEvent{
					Type: types.PodUpdate,
					Pod:  newPod}
			} else {
				log.Println("pod not changed, pod ID is:", newPod.Name)
			}
		}
	}

	return nil
}
