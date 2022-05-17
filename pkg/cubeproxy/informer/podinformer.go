package informer

import (
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
)

type PodInformer interface {
	InitInformer(pods []object.Pod) error
	WatchPodEvent() <-chan types.PodEvent
	InformPod(newPod object.Pod, eType watchobj.EventType) error
	ListPods() []object.Pod
	CloseChan()
}

type ProxyPodInformer struct {
	podChannel chan types.PodEvent
	podCache   map[string]object.Pod

	mtx sync.Mutex
}

func NewPodInformer() PodInformer {
	return &ProxyPodInformer{
		podChannel: make(chan types.PodEvent),
		podCache:   make(map[string]object.Pod),
	}
}

func (i *ProxyPodInformer) InitInformer(pods []object.Pod) error {
	for _, pod := range pods {
		i.podCache[pod.UID] = pod
	}

	return nil
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
	i.mtx.Lock()
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.podChannel <- types.PodEvent{
				Type: types.PodRemove,
				Pod:  newPod,
			}
		} else {
			log.Printf("[INFO]: pod %s not exist, delete do nothing\n", newPod.UID)
		}
	} else {
		// FIXME: If a pod lose its ip..
		// Just handle pod whose ip has been allocated
		if newPod.Status == nil || newPod.Status.IP == nil {
			log.Printf("Pod %v without ip, just ignore", newPod.UID)
			i.mtx.Unlock()
			return nil
		}

		// else cached and judge type
		i.podCache[newPod.UID] = newPod

		if !exist {
			i.podChannel <- types.PodEvent{
				Type: types.PodCreate,
				Pod:  newPod,
			}
		} else {
			// compute pod change: IP / Label
			if object.ComputePodNetworkChange(&newPod, &oldPod) {
				log.Println("[INFO]: pod changed, pod ID is:", newPod.UID)
				i.podChannel <- types.PodEvent{
					Type: types.PodUpdate,
					Pod:  newPod}
			} else {
				log.Println("[INFO]: pod not changed, pod ID is:", newPod.Name)
			}
		}
	}

	i.mtx.Unlock()
	return nil
}
