package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/cubelet/informer/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type PodInformer interface {
	ListAndWatchPodsWithRetry()
	WatchPodEvent() <-chan types.PodEvent
	ListPods() []object.Pod
}

const WatchRetryIntervalSec = 10

func NewPodInformer() (PodInformer, error) {
	return &cubePodInformer{
		podEvent: make(chan types.PodEvent),
		podCache: make(map[string]object.Pod),
	}, nil
}

type cubePodInformer struct {
	nodeUID  string
	podEvent chan types.PodEvent
	podCache map[string]object.Pod
}

func (i *cubePodInformer) ListAndWatchPodsWithRetry() {
	defer close(i.podEvent)
	for {
		i.tryListAndWatchPods()
		time.Sleep(WatchRetryIntervalSec * time.Second)
	}
}

func (i *cubePodInformer) tryListAndWatchPods() {
	// list all pods from apiserver first,
	// in case of cubelet restart or lost connection with apiserver
	// ensure informer cache all pods of apiserver
	if allPods, err := crudobj.GetPods(); err != nil {
		log.Printf("[INFO]: fail to get all pods from apiserver: %v\n", err)
		log.Printf("[INFO]: will retry after %d seconds...\n", WatchRetryIntervalSec)
		return
	} else {
		// put all pods to pod cache
		for _, pod := range allPods {
			if pod.Status != nil && pod.Status.NodeUID == i.nodeUID {
				// simply initialize cache without inform
				// informer could lead to Create event because cache is empty
				// we assume that no new pod will be bound since apiserver lost connection with this Node
				// much simplified :)
				i.podCache[pod.UID] = pod
			}
		}
	}

	// then watch pod status change
	ch, cancel, err := watchobj.WatchPods()
	if err != nil {
		log.Printf("fail to watch pods from apiserver: %v\n", err)
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
				log.Println("[INFO]: Pod caught, but status is nil so Cubelet doesn't handle it")
				continue
			}
			if podEvent.EType == watchobj.EVENT_DELETE || podEvent.Pod.Status.NodeUID == i.nodeUID {
				log.Println("[INFO]: my pod caught, types is", podEvent.EType)
				switch podEvent.EType {
				case watchobj.EVENT_PUT, watchobj.EVENT_DELETE:
					err := i.informPod(podEvent.Pod, podEvent.EType)
					if err != nil {
						return
					}
				default:
					log.Panic("[Error]: Unsupported types in watch pod.")
				}
			} else {
				log.Printf("[INFO]: pod caught, but not my pod, pod UUID = %v, my UUID = %v",
					podEvent.Pod.Status.NodeUID, i.nodeUID)
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func (i *cubePodInformer) WatchPodEvent() <-chan types.PodEvent {
	return i.podEvent
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

func (i *cubePodInformer) informPod(newPod object.Pod, eType watchobj.EventType) error {
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.podEvent <- types.PodEvent{
				Type: types.Remove,
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
				Type: types.Create,
				Pod:  newPod}
		} else {
			// compute pod change: Name / Label / Spec
			if object.ComputeObjectMetaChange(&newPod.ObjectMeta, &oldPod.ObjectMeta) ||
				object.ComputePodSpecChange(&newPod.Spec, &oldPod.Spec) {
				log.Printf("pod %s spec configured\n", newPod.Name)
				i.podEvent <- types.PodEvent{
					Type: types.Update,
					Pod:  newPod}
			} else {
				log.Printf("pod %s spec not change\n", newPod.Name)
			}
		}
	}

	return nil
}
