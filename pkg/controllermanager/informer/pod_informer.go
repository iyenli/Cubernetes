package informer

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/watchobj"
	"Cubernetes/pkg/controllermanager/phase"
	"Cubernetes/pkg/controllermanager/types"
	"Cubernetes/pkg/object"
	"log"
	"time"
)

type PodInformer interface {
	ListAndWatchPodsWithRetry()
	WatchPodEvent() <-chan types.PodEvent
	CloseChan(<-chan types.PodEvent)
	SelectPods(selector map[string]string) []object.Pod
	RecordRemove(uid string)
}

const (
	watchPodsRetryIntervalSec = 8
	watchRSRetryIntervalSec   = 10
	watchASRetryIntervalSec   = 12
)

func NewPodInformer() (PodInformer, error) {
	return &cmPodInformer{
		podEventChans: make([]chan types.PodEvent, 0),
		podCache:      make(map[string]object.Pod),
		rmCache:       make(map[string]interface{}),
	}, nil
}

type cmPodInformer struct {
	podEventChans []chan types.PodEvent
	podCache      map[string]object.Pod
	rmCache       map[string]interface{}
}

func (i *cmPodInformer) ListAndWatchPodsWithRetry() {
	for {
		i.tryListAndWatchPods()
		time.Sleep(watchPodsRetryIntervalSec * time.Second)
	}
}

func (i *cmPodInformer) tryListAndWatchPods() {
	// List all pods from apiserver
	if allPods, err := crudobj.GetPods(); err != nil {
		log.Printf("[Manager] fail to get all pods from apiserver: %v\n", err)
		log.Printf("[Manager] will retry after %d seconds...\n", watchPodsRetryIntervalSec)
		return
	} else {
		for _, pod := range allPods {
			if pod.Status != nil && !phase.NotHandle(pod.Status.Phase) {
				i.podCache[pod.UID] = pod
			}
		}
	}

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
				log.Printf("lost connection with APIServer, retry after %d seconds...\n", watchPodsRetryIntervalSec)
				return
			}
			pod := podEvent.Pod
			log.Printf("manager pod informer get pod %s, event is %s\n", pod.UID, podEvent.EType)
			// pod status not ready to handle by controller_manager
			if (pod.Status == nil || phase.NotHandle(pod.Status.Phase)) &&
				podEvent.EType != watchobj.EVENT_DELETE {
				continue
			}
			if _, ok := i.rmCache[pod.UID]; ok {
				continue
			}
			switch podEvent.EType {
			case watchobj.EVENT_DELETE, watchobj.EVENT_PUT:
				err := i.informPod(pod, podEvent.EType)
				if err != nil {
					log.Println("[INFO]: Delete or put an pod error", pod.UID)
					return
				}
			default:
				log.Panic("Unsupported types in watch pod.")
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (i *cmPodInformer) WatchPodEvent() <-chan types.PodEvent {
	log.Printf("pod informer make a new chan!\n")
	newChan := make(chan types.PodEvent, 10)
	i.podEventChans = append(i.podEventChans, newChan)
	return newChan
}

func (i *cmPodInformer) informPod(newPod object.Pod, eType watchobj.EventType) error {
	oldPod, exist := i.podCache[newPod.UID]

	if eType == watchobj.EVENT_DELETE {
		if exist {
			delete(i.podCache, newPod.UID)
			i.informAll(types.PodEvent{
				Type: types.PodKilled,
				Pod:  oldPod,
			})
		} else {
			log.Printf("pod %s not exist, DELETE do nothing\n", newPod.Name)
		}
	}

	if eType == watchobj.EVENT_PUT {
		if !exist && phase.Running(newPod.Status.Phase) {
			i.podCache[newPod.UID] = newPod
			i.informAll(types.PodEvent{
				Type: types.PodCreate,
				Pod:  newPod,
			})
		} else if exist {
			i.podCache[newPod.UID] = newPod
			log.Printf("Pod %s cpu usage is %v\n", newPod.UID, newPod.Status.ActualResourceUsage.ActualCPUUsage)
			newRunning := phase.Running(newPod.Status.Phase)
			oldRunning := phase.Running(oldPod.Status.Phase)
			if newRunning && oldRunning {
				// both running: update
				i.informAll(types.PodEvent{
					Type: types.PodUpdate,
					Pod:  newPod,
				})
			} else if newRunning && !oldRunning {
				// new running pod
				i.informAll(types.PodEvent{
					Type: types.PodCreate,
					Pod:  newPod,
				})
			} else if !newRunning && oldRunning {
				// pod not running anymore: kill
				delete(i.podCache, newPod.UID)
				i.informAll(types.PodEvent{
					Type: types.PodKilled,
					Pod:  oldPod,
				})
			} // do nothing when both dead
		}
	}

	return nil
}

func (i *cmPodInformer) CloseChan(ch <-chan types.PodEvent) {
	log.Printf("pod informer close a chan!\n")
	found := -1
	for idx, c := range i.podEventChans {
		if c == ch {
			close(c)
			found = idx
			break
		}
	}
	if found != -1 {
		i.podEventChans = append(i.podEventChans[:found], i.podEventChans[found+1:]...)
	}
}

func (i *cmPodInformer) SelectPods(selector map[string]string) []object.Pod {
	matchedPods := make([]object.Pod, 0)
	for _, pod := range i.podCache {
		if object.MatchLabelSelector(selector, pod.Labels) {
			matchedPods = append(matchedPods, pod)
		}
	}
	return matchedPods
}

func (i *cmPodInformer) RecordRemove(uid string) {
	i.rmCache[uid] = true
	delete(i.podCache, uid)
}

func (i *cmPodInformer) informAll(event types.PodEvent) {
	for _, c := range i.podEventChans {
		c <- event
	}
}
