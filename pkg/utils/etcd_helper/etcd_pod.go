package etcd_helper

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
)

const podPrefix = "/Cubernetes/apis/pods/"

func dummy(...any) string { return "" }

func StorePod(ctx *ETCDContext, podName string, pod string) (bool, error) {
	key := podPrefix + podName

	putResponse, err := ctx.Client.KV.Put(context.TODO(), key, pod)
	if err != nil {
		log.Printf("Put api object of pod failed, prev Value: %s \n CreateRevision : %d \n ModRevision: %d \n Version: %d \n",
			string(putResponse.PrevKv.Value), putResponse.PrevKv.CreateRevision, putResponse.PrevKv.ModRevision, putResponse.PrevKv.Version)
		return false, err
	}

	dummy(key, putResponse)
	return true, nil
}

/* podName can be accurate or prefix */
func GetPods(ctx *ETCDContext, podName string) ([][]byte, error) {
	key := podPrefix + podName

	getResponse, err := ctx.Client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		log.Printf("Get api object of pod failed, Key: %s \n", key)
		return nil, err
	}

	var res [][]byte
	for _, resp := range getResponse.Kvs {
		res = append(res, resp.Value)
	}
	return res, nil
}

func GetAllPods(ctx *ETCDContext) ([][]byte, error) {
	return GetPods(ctx, "")
}

func GetPodsRange(ctx *ETCDContext, podNameStart string, podNameEnd string) ([][]byte, error) {
	getResponse, err := ctx.Client.Get(context.TODO(), podPrefix+podNameStart, clientv3.WithRange(podPrefix+podNameEnd))
	if err != nil {
		log.Printf("Get api object of pod in range failed, Key: %s - %s \n", podNameStart, podNameEnd)
		return nil, err
	}

	var res [][]byte
	for _, resp := range getResponse.Kvs {
		res = append(res, resp.Value)
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
