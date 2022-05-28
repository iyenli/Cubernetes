package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubeproxy/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"sync"
	"time"
)

type PodInformer interface {
	WatchPodEvent() <-chan types.PodEvent
	ListAndWatchPodsWithRetry()
	ListPods() []object.Pod
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

func (i *ProxyPodInformer) ListAndWatchPodsWithRetry() {
	defer close(i.podChannel)
	for {
		i.tryListAndWatchPods()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (i *ProxyPodInformer) tryListAndWatchPods() {
	if allPods, err := crudobj.GetPods(); err != nil {
		log.Printf("[INFO]: fail to get all pods from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		for _, pod := range allPods {
			if pod.Status != nil {
				i.podCache[pod.UID] = pod
			}
		}
	}

	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Println("[Error]: Error occurs when watching pods")
		return
	}
	defer cancel()

	for {
		select {
		case podEvent, ok := <-ch:
			if !ok {
				log.Printf("[INFO]: lost connection with APIServer, retry after %d seconds...\n", WatchRetryIntervalSec)
				return
			}
			if podEvent.Pod.Status == nil && podEvent.EType != watchobj.EVENT_DELETE {
				log.Println("[INFO]: Pod caught, but status is nil so Cubeproxy doesn't handle it")
				continue
			}
			switch podEvent.EType {
			case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
				err := i.informPod(podEvent.Pod, podEvent.EType)
				if err != nil {
					log.Println("[Error]: Error when inform pod: ", podEvent.Pod.UID)
					return
				}
			default:
				log.Panic("[Fatal]: Unsupported types in watch pod")
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *ProxyPodInformer) WatchPodEvent() <-chan types.PodEvent {
	return i.podChannel
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

func (i *ProxyPodInformer) informPod(newPod object.Pod, eType watchobj.EventType) error {
	i.mtx.Lock()
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.podChannel <- types.PodEvent{
				Type: types.Remove,
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
				Type: types.Create,
				Pod:  newPod,
			}
		} else {
			// compute pod change: IP / Label
			if object.ComputePodNetworkChange(&newPod, &oldPod) {
				log.Println("[INFO]: pod changed, pod ID is:", newPod.UID)
				i.podChannel <- types.PodEvent{
					Type: types.Update,
					Pod:  newPod}
			} else {
				log.Println("[INFO]: pod not changed, pod ID is:", newPod.Name)
			}
		}
	}

	i.mtx.Unlock()
	return nil
}
