package etcd_helper

import (
	"Cubernetes/pkg/object"
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"log"
)

const podPrefix = "/Cubernetes/apis/pods/"

func dummy(...any) string { return "" }

func StorePod(ctx *ETCDContext, podName string, pod object.Pod) (bool, error) {
	key := podPrefix + podName

	podJson, err := json.Marshal(pod)
	if err != nil {
		log.Println("Json marshall failed")
		return false, err
	}

	putResponse, err := ctx.Client.KV.Put(context.TODO(), key, string(podJson))
	if err != nil {
		log.Printf("Put api object of pod failed, prev Value: %s \n CreateRevision : %d \n ModRevision: %d \n Version: %d \n",
			string(putResponse.PrevKv.Value), putResponse.PrevKv.CreateRevision, putResponse.PrevKv.ModRevision, putResponse.PrevKv.Version)
		return false, err
	}

	return true, nil
}

// GetPods podName can be accurate or prefix
func GetPods(ctx *ETCDContext, podName string) ([]object.Pod, error) {
	key := podPrefix + podName

	getResponse, err := ctx.Client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		log.Printf("Get api object of pod failed, Key: %s \n", key)
		return nil, err
	}

	var res []object.Pod
	for _, resp := range getResponse.Kvs {
		res = append(res, object.Pod{})
		err := json.Unmarshal(resp.Value, &res[len(res)-1])
		if err != nil {
			log.Printf("Pod unmarshal failed. error: %s \n, value: %s", err, resp.Value)
			return nil, err
		}
	}
	return res, nil
}

func GetAllPods(ctx *ETCDContext) ([]object.Pod, error) {
	return GetPods(ctx, "")
}

func GetPodsRange(ctx *ETCDContext, podNameStart string, podNameEnd string) ([]object.Pod, error) {
	getResponse, err := ctx.Client.Get(context.TODO(), podPrefix+podNameStart, clientv3.WithRange(podPrefix+podNameEnd))
	if err != nil {
		log.Printf("Get api object of pod in range failed, Key: %s - %s \n", podNameStart, podNameEnd)
		return nil, err
	}

	var res []object.Pod
	for _, resp := range getResponse.Kvs {
		res = append(res, object.Pod{})
		err := json.Unmarshal(resp.Value, &res[len(res)-1])
		if err != nil {
			log.Printf("Pod unmarshal failed.\n")
			return nil, err
		}
	}
	return res, nil
}

func DeletePod(ctx *ETCDContext, podName string) (bool, error) {
	key := podPrefix + podName
	deleteResp, err := ctx.Client.KV.Delete(context.TODO(), key)
	if err != nil {
		log.Printf("Delete api object of pod failed, Key: %s \n", key)
		return false, err
	}

	return deleteResp.Deleted > 0, nil
}

func GetPodWatcher(ctx *ETCDContext, podName string, accurate bool) (res clientv3.WatchChan) {
	key := podPrefix + podName
	if accurate {
		res = ctx.Client.Watcher.Watch(context.TODO(), key)
	} else {
		res = ctx.Client.Watcher.Watch(context.TODO(), key, clientv3.WithPrefix())
	}
	return
}
